// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package chash

// A custom error type called `err`. It has two fields: Msg (a string
// message) and Code (an integer code).
type err struct {
	Msg  string
	Code int
}

// String and Error. Both return the
// error message stored in the `Msg` field.
func (e err) String() string {
	return e.Msg
}

// String and Error. Both return the
// error message stored in the `Msg` field.
func (e err) Error() string {
	return e.Msg
}

// Several global variables that represent common errors that may be
// returned by the CHash functions.
var (
	ErrGroupNotFound   = err{Code: 10000, Msg: "group not found"}
	ErrGroupExisted    = err{Code: 10001, Msg: "group already existed"}
	ErrNoResultMatched = err{Code: 10002, Msg: "no result matched"}
	ErrKeyExisted      = err{Code: 10003, Msg: "key already existed"}
)
