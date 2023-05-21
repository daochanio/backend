package usecases

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/common"
)

type HydrateUsers struct {
	logger     common.Logger
	blockchain Blockchain
	database   Database
	storage    Storage
	proxy      SafeProxy
}

type HydrateUsersInput struct {
	Addresses []string
}

func NewHydrateUsersUseCase(logger common.Logger, blockchain Blockchain, database Database, storage Storage, proxy SafeProxy) *HydrateUsers {
	return &HydrateUsers{
		logger,
		blockchain,
		database,
		storage,
		proxy,
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
//
// 1. Fetch the avatar text record from ENS using the name
//
// 2. If nft uri detected, parse information and fetch the nft metadata uri from the contract
// then follow the metadata uri to get the image url.
//
// 3. The resulting image url is hashed to derive a unique id
//
// 4. If IPFS scheme detected on the url, resolve the https ipfs url
//
// 5. If the unique id is different from the current id on the user
// the image url is then uploaded to the CDN and the user record is updated with the lates url
func (u *HydrateUsers) Execute(ctx context.Context, input HydrateUsersInput) {
	// We dedupe addresses to ensure we only processes each address once regardless of multiple updates
	addresses := map[string]bool{
		"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045": true,
	}
	for _, address := range input.Addresses {
		addresses[address] = true
	}

	u.logger.Info(ctx).Msgf("hydrating %v users", len(addresses))

	for address := range addresses {
		name, err := u.hydrateName(ctx, address)

		if err != nil {
			u.logger.Warn(ctx).Err(err).Msgf("skipping name hydration for %v", address)
			continue
		}

		avatar, err := u.hydrateAvatar(ctx, name)

		if err != nil {
			u.logger.Warn(ctx).Err(err).Msgf("skipping avatar hydration for %v", address)
			continue
		}

		if err = u.database.UpdateUser(ctx, address, name, avatar); err != nil {
			u.logger.Error(ctx).Err(err).Msgf("error hydrating user %v", address)
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

	// TODO: only supporting mainnet right now (i.e chain id 1)
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

		nftImageURI, err := u.proxy.GetNFTImageURI(ctx, nftURI)

		if err != nil {
			return nil, err
		}

		uri = &nftImageURI
	}

	if err != nil {
		return nil, err
	}

	data, contentType, err := u.proxy.GetData(ctx, *uri)

	if err != nil {
		return nil, err
	}

	fileName := u.getFileName(*uri, contentType)

	image, err := u.storage.UploadImage(ctx, fileName, contentType, data)

	if err != nil {
		return nil, err
	}

	u.logger.Info(ctx).Msgf("uploaded avatar: %s for name: %s from uri: %s", image.Url(), *name, *uri)

	return &image, nil
}

// hash the uri to derive a unique but idempotent filename
func (h *HydrateUsers) getFileName(uri string, contentType string) string {
	hash := sha256.New()
	hash.Write([]byte(uri))
	name := hex.EncodeToString(hash.Sum(nil))
	ext := strings.Split(contentType, "/")[1]
	return fmt.Sprintf("%v.%v", name, ext)
}
