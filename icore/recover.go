package icore

import "blockchain/models"

type Recover interface {
	RecvReq(to models.Address, req models.RecoverReq) error
	RecvVerifyReq(to models.Address, req models.RecoverVerifyReq) error
	RecvRes(to models.Address, req models.RecoverRes) error
	RecvMd5Req(to models.Address, req models.RecoverMd5Req) error
	RecvMd5Res(to models.Address, req models.RecoverMd5Res) error
	VerifyMd5Req(to models.Address, req models.VerifyMd5Req) error
	VerifyMd5Res(to models.Address, req models.VerifyMd5Res) error
}
