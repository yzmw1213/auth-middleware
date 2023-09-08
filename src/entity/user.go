package entity

type User struct {
	UserID      int64  `json:"user_id"`      // ユーザーID
	Name        string `json:"name"`         //名前
	Email       string `json:"email"`        //Eメール
	FirebaseUID string `json:"firebase_uid"` // firebase id
	DeleteFlag  bool   `json:"delete_flag"`  // 削除フラグ
}
