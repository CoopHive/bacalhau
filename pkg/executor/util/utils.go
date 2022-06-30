package util

import (
	"github.com/filecoin-project/bacalhau/pkg/executor"
	"github.com/filecoin-project/bacalhau/pkg/executor/docker"
	"github.com/filecoin-project/bacalhau/pkg/executor/language"
	noop_executor "github.com/filecoin-project/bacalhau/pkg/executor/noop"
	"github.com/filecoin-project/bacalhau/pkg/executor/python_wasm"
	"github.com/filecoin-project/bacalhau/pkg/storage"
	"github.com/filecoin-project/bacalhau/pkg/storage/ipfs/apicopy"
	"github.com/filecoin-project/bacalhau/pkg/system"
)

func NewStandardExecutors(
	cm *system.CleanupManager,
	ipfsMultiAddress,
	dockerID string,
) (map[executor.EngineType]executor.Executor, error) {

	// Don't allow user to choose the fuse driver in case it has security issues.
	// ipfsFuseStorage, err := fusedocker.NewStorageProvider(cm, ipfsMultiAddress)
	// if err != nil {
	// 	return nil, err
	// }

	ipfsAPICopyStorage, err := apicopy.NewStorageProvider(cm, ipfsMultiAddress)
	if err != nil {
		return nil, err
	}

	exDocker, err := docker.NewExecutor(cm, dockerID,
		map[string]storage.StorageProvider{
			// fuse driver is disabled so that - in case it poses a security
			// risk - arbitrary users can't request it
			// storage.IPFS_FUSE_DOCKER: ipfsFuseStorage,
			storage.IPFSAPICopy: ipfsAPICopyStorage,
			// we make the copy driver the "default" storage driver for docker
			// users have to specify the fuse driver explicitly
			storage.IPFSDefault: ipfsAPICopyStorage,
		})
	if err != nil {
		return nil, err
	}

	executors := map[executor.EngineType]executor.Executor{
		executor.EngineDocker: exDocker,
	}

	// language executors wrap other executors, so pass them a reference to all
	// the executors so they can look up the ones they need
	exLang, err := language.NewExecutor(cm, executors)
	executors[executor.EngineLanguage] = exLang
	if err != nil {
		return nil, err
	}
	exPythonWasm, err := python_wasm.NewExecutor(cm, executors)
	executors[executor.EnginePythonWasm] = exPythonWasm
	if err != nil {
		return nil, err
	}
	return executors, nil
}

// return noop executors for all engines
func NewNoopExecutors(
	cm *system.CleanupManager,
) (map[executor.EngineType]executor.Executor, error) {
	noopExecutor, err := noop_executor.NewExecutor()

	if err != nil {
		return nil, err
	}

	return map[executor.EngineType]executor.Executor{
		executor.EngineDocker: noopExecutor,
		executor.EngineNoop:   noopExecutor,
	}, nil
}
