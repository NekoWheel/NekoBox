// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"
	"testing"
)

func TestQuestions(t *testing.T) {
	t.Parallel()

	db, cleanup := newTestDB(t)
	ctx := context.Background()

	usersStore := NewUsersStore(db)

	for _, tc := range []struct {
		name string
		test func(t *testing.T, ctx context.Context, db *users)
	}{
		//Create(ctx context.Context, opts CreateQuestionOptions) (*Question, error)
		//GetByID(ctx context.Context, id uint) (*Question, error)
		//GetByUserID(ctx context.Context, userID uint, opts GetQuestionsByUserIDOptions) ([]*Question, error)
		//GetByAskUserID(ctx context.Context, userID uint, opts GetQuestionsByAskUserIDOptions) ([]*Question, error)
		//AnswerByID(ctx context.Context, id uint, answer string) error
		//DeleteByID(ctx context.Context, id uint) error
		//UpdateCensor(ctx context.Context, id uint, opts UpdateQuestionCensorOptions) error
		//Count(ctx context.Context, userID uint, opts GetQuestionsCountOptions) (int64, error)

		//{"Create", testUsersCreate},
		//{"GetByID", testUsersGetByID},
		//{"GetByEmail", testUsersGetByEmail},
		//{"GetByDomain", testUsersGetByDomain},
		//{"Update", testUsersUpdate},
		//{"UpdateHarassmentSetting", testUsersUpdateHarassmentSetting},
		//{"Authenticate", testUsersAuthenticate},
		//{"ChangePassword", testUsersChangePassword},
		//{"UpdatePassword", testUsersUpdatePassword},
		//{"Deactivate", testUsersDeactivate},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				if err := cleanup("users"); err != nil {
					t.Fatal(err)
				}
			})
			tc.test(t, ctx, usersStore.(*users))
		})
	}
}
