package saga

import (
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	error2 "service/internal/pkg/error"
	"service/internal/pkg/grpc/example"
	"service/internal/pkg/saga/grpc"
	"service/internal/pkg/saga/privateapi"
)

type TestingSaga struct {
	testingGRPC *grpc.TestingGRPC
	testingAPI  privateapi.TestingAPI

	testingGRPCCleanup func()

	testingRPCRollBack []byte
	testingAPIRollBack interface{}
}

/** --- NEW CLIENT --- */

func (saga *TestingSaga) NewTestingClient() {
	saga.testingGRPC, saga.testingGRPCCleanup = grpc.NewTestingGRPC()
	saga.testingAPI = privateapi.NewTestingAPI()
}

/** --- ITEM SERVICE CLIENT --- */

func (saga *TestingSaga) TestingStore(request *example.TestingRequest) (string, []byte) {
	check, err := saga.testingGRPC.Testing.Store(saga.testingGRPC.Ctx, request)
	if err != nil {
		saga.testingGRPC = nil
		saga.testingGRPCCleanup()

		error2.ErrXtremeTestingSave(err.Error())
	}

	saga.testingRPCRollBack = check.Result

	return check.Message, check.Result
}

func (saga *TestingSaga) TestingStoreAPI(request *example.TestingRequest) interface{} {
	resp := saga.testingAPI.Store(request)
	result := resp.Result
	if result != nil {
		saga.testingAPIRollBack = result
	}

	return result
}

func (saga *TestingSaga) TestingRollbackStore() {
	_, err := saga.testingGRPC.Testing.RollbackStore(saga.testingGRPC.Ctx, &example.RollBackRequest{Data: saga.testingRPCRollBack})
	if err != nil {
		xtremepkg.LogError(err, true)
	}
}

func (saga *TestingSaga) TestingAPIRollbackStore() {
	saga.testingAPI.RollBack(saga.testingAPIRollBack)
}

/** --- DEFER FUNCTION --- */

func (saga *TestingSaga) Close() {
	if r := recover(); r != nil {
		if saga.testingGRPC != nil && len(saga.testingRPCRollBack) > 0 {
			saga.TestingRollbackStore()
			saga.testingGRPCCleanup()
		}

		if saga.testingAPIRollBack != nil {
			saga.TestingAPIRollbackStore()
		}

		panic(r)
	}

	saga.Cleanup()
}

func (saga *TestingSaga) Cleanup() {
	if saga.testingGRPC != nil {
		saga.testingGRPCCleanup()
	}
}
