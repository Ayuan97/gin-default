package models

import (
	"encoding/base64"
	"encoding/json"
	"gorm.io/gorm"
	"justus/pkg/aes"
	"justus/pkg/setting"
	"strconv"
	"strings"
)

type PictureLibrary struct {
	Model
	ImgUrl     string         `json:"img_url" gorm:"column:img_url"`
	SubImgUrl  string         `json:"sub_img_url"`
	Type       int            `json:"type" gorm:"column:type"`
	Uid        int            `json:"uid" gorm:"column:uid"`
	HotNum     int            `json:"hot_num" gorm:"column:hot_num"`
	InTuneType string         `json:"in_tune_type"`
	Position   string         `json:"position"`
	IsVisible  int            `json:"is_visible"`
	VideoUrl   string         `json:"video_url"`
	Content    string         `json:"content"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
}
type PictureLibraryList struct {
	ID            int                `json:"id"`
	ImgUrl        string             `json:"img_url"`
	SubImgUrl     string             `json:"sub_img_url"`
	Type          int                `json:"type"`
	Uid           int                `json:"uid"`
	FirstName     string             `json:"first_name"`
	LastName      string             `json:"last_name"`
	Avatar        string             `json:"avatar"`
	IsLike        int                `json:"is_like"`
	IsFollow      int                `json:"is_follow"`
	IsFriend      int                `json:"is_friend"`
	IsFuzzy       int                `json:"is_fuzzy"`
	Source        int                `json:"source"`
	IsTuneSuccess int                `json:"is_tune_success"`
	Comment       []*CommentFormated `json:"comment"`
	Like          []*LikeFormated    `json:"like"`
	Topics        []TopicPicture     `json:"topics"`
	HotNum        int                `json:"hot_num"`
	InTuneType    string             `json:"in_tune_type"`
	Position      string             `json:"position"`
	IsVisible     int                `json:"is_visible"`
	VideoUrl      string             `json:"video_url"`
	LikeNum       int                `json:"like_num"`
	CommentNum    int                `json:"comment_num"`
	Content       string             `json:"content"`
	ImgCode       string             `json:"img_code"`
	SubImgCode    string             `json:"sub_img_code"`
	CreatedAt     int                `json:"created_at"`
}

func (p *PictureLibrary) getUrl() string {
	if p.ImgUrl == "" {
		return ""
	}
	if strings.Contains(p.ImgUrl, "http") {
		return p.ImgUrl
	} else {
		return setting.AppSetting.ImageUrl + "/" + p.ImgUrl
	}
}
func (p *PictureLibrary) getSubUrl() string {
	if p.SubImgUrl == "" {
		return ""
	}
	if strings.Contains(p.SubImgUrl, "http") {
		return p.SubImgUrl
	} else {
		return setting.AppSetting.ImageUrl + "/" + p.SubImgUrl
	}
}
func (p *PictureLibrary) getVideoUrl() string {
	if p.VideoUrl == "" {
		return ""
	}
	if strings.Contains(p.VideoUrl, "http") {
		return p.VideoUrl
	} else {
		return setting.AppSetting.ImageUrl + "/" + p.VideoUrl
	}
}

func GetImgCode(t int, PId int, p int) string {
	if t == 1 {
		imgInfo := make(map[string]string)
		imgInfo["id"] = strconv.Itoa(PId)
		imgInfo["type"] = strconv.Itoa(p)
		jsonStr, _ := json.Marshal(imgInfo)
		str := aes.AesEncryptCBC(jsonStr)
		encrypt := base64.StdEncoding.EncodeToString(str)
		return encrypt
	} else {
		return ""
	}
}

func (p *PictureLibrary) Format() *PictureLibraryList {
	if p.Uid <= 0 {
		return nil
	}

	return &PictureLibraryList{
		ID:         p.ID,
		ImgUrl:     p.getUrl(),
		SubImgUrl:  p.getSubUrl(),
		Type:       p.Type,
		Uid:        p.Uid,
		HotNum:     p.HotNum,
		InTuneType: p.InTuneType,
		Position:   p.Position,
		IsVisible:  p.IsVisible,
		VideoUrl:   p.getVideoUrl(),
		Content:    p.Content,
		ImgCode:    GetImgCode(p.Type, p.ID, 1),
		SubImgCode: GetImgCode(p.Type, p.ID, 2),
		CreatedAt:  p.CreatedAt,
	}
}

// GetArticle Get a single article based on ID
func GetPictureLibrary(id int) (*PictureLibrary, error) {
	var pictureLibrary PictureLibrary
	err := db.Where("id = ?", id).First(&pictureLibrary).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &pictureLibrary, nil
}

// GetPictureOne Get a single article based on ID
func (p *PictureLibrary) GetPictureOne() (*PictureLibrary, error) {
	var pictureLibrary PictureLibrary
	err := db.Where("id = ?", p.ID).First(&pictureLibrary).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &pictureLibrary, nil
}

// GetImgList GetArticleList Get a list of articles
func (p *PictureLibrary) GetImgList(offset int, limit int) ([]*PictureLibrary, error) {
	var pictureLibrary []*PictureLibrary
	err := db.Where("uid = ?", p.Uid).Where("is_visible = 1").Order("id desc").Offset(offset).Limit(limit).Find(&pictureLibrary).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return pictureLibrary, nil
}

func (p *PictureLibrary) GetPictureLibrarys(pid []string, offset int, limit int) ([]*PictureLibrary, error) {
	var pictureLibrary []*PictureLibrary
	err := db.Where("id in (?)", pid).Offset(offset).Limit(limit).Find(&pictureLibrary).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return pictureLibrary, nil
}

// 通过图片ID 查询话题信息
func (m *PictureLibrary) GetPictureKeyListByIds(pid []int) map[int]*PictureLibrary {
	var pictureLibrary []PictureLibrary
	result := make(map[int]*PictureLibrary)
	if len(pid) > 0 {
		err := db.Where("id in (?)", pid).Order("id desc").Find(&pictureLibrary).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return result
		}
		for _, v := range pictureLibrary {
			temp := v
			result[v.ID] = &temp
		}
	}
	return result
}
