package deploylib

import (
	"homedy/config"
	"homedy/internal/libs/cryptolib"
	"homedy/internal/models"

	"github.com/compose-spec/compose-go/v2/types"
)

func cryptoEnv(env types.MappingWithEquals, f func([]byte, []byte) (string, error)) (types.MappingWithEquals, error) {
	cryptoMap := make(types.MappingWithEquals, len(env))
	for k, v := range env {
		cryptoVal, err := f([]byte(*v), []byte(config.DEPLOY_ENV_CRYPTO_KEY))
		if err != nil {
			return nil, err
		}
		cryptoMap[k] = &cryptoVal
	}
	return cryptoMap, nil
}

func EncryptEnv(env types.MappingWithEquals) (types.MappingWithEquals, error) {
	return cryptoEnv(env, cryptolib.EncryptGCM)
}

func DecryptEnv(env types.MappingWithEquals) (types.MappingWithEquals, error) {
	return cryptoEnv(env, cryptolib.DecryptGCM)
}

func cryptoSessionEnv(env models.DeploySessionEnv, f func(types.MappingWithEquals) (types.MappingWithEquals, error)) (*models.DeploySessionEnv, error) {
	cryptoGlobal, err := f(types.MappingWithEquals(env.Global))
	if err != nil {
		return nil, err
	}

	cryptoRepo := make(models.RepoEnv, len(env.Repo))
	for repoID, services := range env.Repo {
		cryptoRepo[repoID] = models.ServiceEnv{}
		for k, v := range services {
			cryptoRepoService, err := f(v)
			if err != nil {
				return nil, err
			}
			cryptoRepo[repoID][k] = cryptoRepoService
		}
	}

	return &models.DeploySessionEnv{Global: models.GlobalEnv(cryptoGlobal), Repo: cryptoRepo}, nil
}

func EncryptSessionEnv(env models.DeploySessionEnv) (*models.DeploySessionEnv, error) {
	return cryptoSessionEnv(env, EncryptEnv)
}

func DecryptSessionEnv(env models.DeploySessionEnv) (*models.DeploySessionEnv, error) {
	return cryptoSessionEnv(env, DecryptEnv)
}
