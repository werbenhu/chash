// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package chash

type err struct {
	Msg  string
	Code int
}

func (e err) String() string {
	return e.Msg
}

func (e err) Error() string {
	return e.Msg
}

var (
	ErrGroupNotFound   = err{Code: 10000, Msg: "group not found"}
	ErrGroupExisted    = err{Code: 10001, Msg: "group already existed"}
	ErrNoResultMatched = err{Code: 10002, Msg: "no result matched"}
	ErrKeyExisted      = err{Code: 10003, Msg: "key already existed"}
)
