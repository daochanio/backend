package usecases

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
)

type HydrateUsers struct {
	logger     common.Logger
	blockchain Blockchain
	database   Database
	storage    Storage
	safeProxy  SafeProxy
}

type HydrateUsersInput struct {
	Addresses []string
}

func NewHydrateUsersUseCase(logger common.Logger, blockchain Blockchain, database Database, storage Storage, safeProxy SafeProxy) *HydrateUsers {
	return &HydrateUsers{
		logger,
		blockchain,
		database,
		storage,
		safeProxy,
	}
}

// There be dragons below.
//
// The general reasoning behind doing all this work async in the backend is that avatars are
// very dynamic, slow, unreliable and potentially very big.
// The goal is to have a fast, reliable and small profile pictures that can be rendered in a browser very quickly to not
// destroy page load times, SEO and user experience.
// The avatar image is thus resolved, uploaded and cached in our CDN and the url is stored in the database.
// But resolving the actual image url from an avatar text record is quite an involved process.
// See: https://docs.ens.domains/ens-improvement-proposals/ensip-12-avatar-text-records
//
// Steps:
//  1. Fetch the avatar text record from ENS using the name
//  2. This record can point to any arbitrary server, so we proxy all the following http requests through a "safe" proxy to avoid leaking sensitive server information
//  3. If nft uri detected, parse information and fetch the nft metadata uri from the contract then follow the metadata uri to get the image url.
//  4. The resulting image url is hashed to derive a unique and idempotent file name
//  5. If IPFS scheme detected on the url, resolve the https ipfs url
//  6. Check if the file already exists in storage
//  7. If not, download the image and upload it to our storage
//
// TODO:
//   - Supported Data URIs
//   - Support other chains for NFT URIs
func (u *HydrateUsers) Execute(ctx context.Context, input HydrateUsersInput) {
	// We dedupe addresses to ensure we only processes each address once regardless of multiple updates
	addresses := map[string]bool{}
	for _, address := range input.Addresses {
		addresses[address] = true
	}

	u.logger.Info(ctx).Msgf("hydrating %v users", len(addresses))

	for address := range addresses {
		name, err := u.hydrateName(ctx, address)

		if err != nil {
			u.logger.Warn(ctx).Err(err).Msgf("name err - skipping hydration for %v", address)
			continue
		}

		avatar, err := u.hydrateAvatar(ctx, name)

		if err != nil {
			u.logger.Warn(ctx).Err(err).Msgf("avatar err - skipping hydration for %v", address)
			continue
		}

		if err = u.database.UpdateUser(ctx, address, name, avatar); err != nil {
			u.logger.Error(ctx).Err(err).Msgf("error saving user hydration %v", address)
		}
	}
}

func (u *HydrateUsers) hydrateName(ctx context.Context, address string) (*string, error) {
	name, err := u.blockchain.GetNameByAddress(ctx, address)

	if err != nil {
		return nil, err
	}

	return name, nil
}

func (u *HydrateUsers) hydrateAvatar(ctx context.Context, name *string) (*entities.Image, error) {
	// if theres no name, theres also no avatar
	if name == nil {
		return nil, nil
	}

	uri, err := u.blockchain.GetAvatarURIByName(ctx, *name)

	if err != nil {
		return nil, err
	}

	if uri == nil {
		return nil, nil
	}

	if suffix, ok := strings.CutPrefix(*uri, "eip155:1/"); ok {
		standard, info, ok := strings.Cut(suffix, ":")
		if !ok {
			return nil, errors.New("invalid nft uri")
		}

		address, id, ok := strings.Cut(info, "/")
		if !ok {
			return nil, errors.New("invalid nft info")
		}

		nftURI, err := u.blockchain.GetNFTURI(ctx, standard, address, id)

		if err != nil {
			return nil, err
		}

		nftImageURI, err := u.safeProxy.GetNFTImageURI(ctx, nftURI)

		if err != nil {
			return nil, err
		}

		uri = &nftImageURI
	}

	if err != nil {
		return nil, err
	}

	fileName := u.getFileName(*uri)

	if existingImage, err := u.storage.GetImageByFileName(ctx, fileName); err == nil && existingImage != nil {
		u.logger.Info(ctx).Msgf("found existing avatar: %s for name: %s from uri: %s", existingImage.Url(), *name, *uri)
		return existingImage, nil
	}

	data, contentType, err := u.safeProxy.DownloadImage(ctx, *uri)

	if err != nil {
		return nil, err
	}

	image, err := u.storage.UploadImage(ctx, fileName, contentType, data)

	if err != nil {
		return nil, err
	}

	u.logger.Info(ctx).Msgf("uploaded avatar: %s for name: %s from uri: %s", image.Url(), *name, *uri)

	return &image, nil
}

// hash the uri to derive a unique but idempotent filename
func (h *HydrateUsers) getFileName(uri string) string {
	hash := sha256.New()
	hash.Write([]byte(uri))
	return hex.EncodeToString(hash.Sum(nil))
}
