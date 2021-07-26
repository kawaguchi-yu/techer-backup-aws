package infra

import (
	"encoding/json"
	"fmt"
	"hello/server/domain"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func DBInit() error {
	env := getEnv() //envに環境変数を代入
	dsn := env.userName + ":" + env.password + "@tcp(" + env.host + ")/" + env.dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	fmt.Printf("%v\n\n", dsn)
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) //gorm.Openでdbと接続している
	if err != nil {
		return echo.ErrInternalServerError
	}
	db = gormDB
	fmt.Println(db) //コンソールに出力
	// migrate
	db.AutoMigrate(&domain.Tag{}, domain.Good{}, &domain.ExternaiURL{})
	db.AutoMigrate(&domain.Post{})
	db.AutoMigrate(&domain.Profile{})
	db.AutoMigrate(&domain.User{})
	return nil
}
func GetDB() *gorm.DB {
	return db
}

type env struct {
	userName string
	password string
	host     string
	dbName   string
}

func getEnv() env {
	e := env{
		userName: "root",
		password: "hoge",
		host:     "db:3306",
		dbName:   "db",
	}
	return e
}
func DBCreateUser(c echo.Context, db *gorm.DB) error { //渡された値をDBに入れる
	u := new(domain.User)
	c.Bind(u) //cの中のユーザー情報をuに入れる
	rawPassword := []byte(u.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(rawPassword, 4)
	if err != nil {
		return echo.ErrBadRequest
	}
	u.Password = string(hashedPassword)
	fmt.Printf("%+v\n\n", u)
	result := db.Create(&u)
	if result.Error != nil {
		return c.JSON(http.StatusBadRequest, "メールアドレスが重複しています")
	}
	return c.JSON(http.StatusOK, "name:"+u.Name+", email:"+u.EMail+", password:"+u.Password)
}

func Login(c echo.Context, db *gorm.DB) error { //emailとpasswordでjwt入りcookie貰える
	u := new(domain.User)
	c.Bind(u)
	//ここにメルアドがdbにあるかをチェックする処理を書く
	dbPassword, err := getUser(u.EMail, db)
	if err != nil { //errの中身がnil以外なら終わる
		return c.JSON(http.StatusBadRequest, "メールアドレスが存在しませんでした")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword.Password), []byte(u.Password)); err != nil {
		return c.JSON(http.StatusBadRequest, "パスワードが違います")
	}
	JWTToken, err := CreateJWT(u.EMail)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "JWTが生成できませんでした")
	}
	cookie := CreateCookie(JWTToken)
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, cookie)

}

func ReadCookieReturnUser(c echo.Context, db *gorm.DB) error {
	email, err := ReadCookie(c)
	if err != nil {
		fmt.Printf("クッキー読み取りに失敗しました")
		return c.JSON(http.StatusBadRequest, nil)
	}
	var user domain.User
	if err := db.First(&user, "e_mail=?", email).Error; err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	return c.JSON(http.StatusOK, user)
}
func ReadCookieReturnIcon(c echo.Context, db *gorm.DB) error {
	email, err := ReadCookie(c)
	if err != nil {
		fmt.Printf("クッキー読み取りに失敗しました")
		return c.JSON(http.StatusBadRequest, nil)
	}
	var user domain.User
	if err := db.First(&user, "e_mail=?", email).Error; err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	prevDir, err := filepath.Abs(".") //カレントディレクトリのパスを保存
	if err != nil {
		return c.JSON(http.StatusBadRequest, "フォルダパスを取得できませんでした")
	}
	defer os.Chdir(prevDir) //もとに戻る
	err = os.Chdir("img")   //imgディレクトリに移動する
	if err != nil {
		return c.JSON(http.StatusBadRequest, "フォルダを移動できませんでした")
	}
	if user.Icon == "" {
		return c.File("firsticon.jpg")
	}
	fmt.Printf("iconは正常に読み取れました\n")
	return c.File(user.Icon)
}
func SetIcon(c echo.Context, db *gorm.DB) error {
	icon, err := c.FormFile("file") //cからファイルを取り出し
	if err != nil {
		fmt.Printf("ファイルが読み込めません\n")
		return c.JSON(http.StatusBadRequest, icon)
	}
	src, err := icon.Open() //io.Readerに変換
	if err != nil {
		fmt.Printf("ファイルをioに変換できませんでした\n")
		return c.JSON(http.StatusBadRequest, "ファイルをioに変換できませんでした")
	}
	defer src.Close()
	prevDir, err := filepath.Abs(".") //カレントディレクトリのパスを保存
	if err != nil {
		return c.JSON(http.StatusBadRequest, "フォルダパスを取得できませんでした")
	}
	defer os.Chdir(prevDir) //もとに戻る
	err = os.Chdir("img")   //imgディレクトリに移動する
	if err != nil {
		return c.JSON(http.StatusBadRequest, "フォルダを移動できませんでした")
	}
	iconModel := strings.Split(icon.Filename, ".")
	iconName := iconModel[0]
	extension := iconModel[1]
	dst, err := os.Create(fmt.Sprintf("%s_out.%s", iconName, extension)) //新しいファイルを作り、名前を決める
	if err != nil {
		fmt.Printf("ファイルが作れませんでした\n")
		return c.JSON(http.StatusBadRequest, "ファイルが作れませんでした")
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil { //ファイルの内容をコピー
		fmt.Printf("コピーできませんでした\n")
		return c.JSON(http.StatusBadRequest, "コピーできませんでした")
	}

	//ここまでが画像をローカルフォルダに保存する行程、ここからがuserのiconに画像データを入れる
	email, err := ReadCookie(c)
	if err != nil {
		fmt.Printf("クッキー読み取りに失敗しました")
		return c.JSON(http.StatusBadRequest, nil)
	}
	var user domain.User
	if err := db.First(&user, "e_mail=?", email).Error; err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	db.Model(&user).Update("icon", dst.Name())
	fmt.Printf("ユーザーネーム=%v\n", user.Name)
	fmt.Printf("正常に終了しました\n" + dst.Name())
	return c.File(user.Icon)
	//例c.File(test.jpg)→test.jpgのファイルが送られる。
}

func getUser(email string, db *gorm.DB) (domain.User, error) {
	var user domain.User
	if err := db.First(&user, "e_mail = ?", email).Error; err != nil {
		fmt.Printf("メールアドレスが存在しませんでした\n")
		return user, echo.ErrBadRequest
	}
	fmt.Printf("user=%v\n", user)
	return user, nil
}
func GetUserModel(b io.ReadCloser) (domain.User, error) {
	var jsonData = make(map[string]string) //空っぽのmapを作る
	var user domain.User
	//デコードしてio.Reader型に変換する
	if err := json.NewDecoder(b).Decode(&jsonData); err != nil {
		return user, echo.ErrBadRequest
	}
	if jsonData == nil {
		return user, echo.ErrInternalServerError
	}

	name := jsonData["Name"]
	eMail := jsonData["EMail"]
	rawPassword := []byte(jsonData["Password"])
	//bcryptでハッシュ化したパスワードをhashedPasswordに入れる
	hashedPassword, err := bcrypt.GenerateFromPassword(rawPassword, 4)
	if err != nil {
		return user, echo.ErrBadRequest
	}

	user.Name = name
	user.EMail = eMail
	user.Password = string(hashedPassword)
	return user, nil
}
