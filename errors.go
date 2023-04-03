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
	ErrBucketNotFound  = err{Code: 10000, Msg: "bucket not found"}
	ErrBucketExisted   = err{Code: 10001, Msg: "bucket already existed"}
	ErrNoResultMatched = err{Code: 10002, Msg: "no result matched"}
)
