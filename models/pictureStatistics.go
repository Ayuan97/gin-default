package models

import (
	"fmt"
	"gorm.io/gorm"
	"justus/global"
)

type PictureStatistics struct {
	Model
	PId            int    `json:"p_id"`
	ImpressionNum            int    `json:"impression_num"`
	CurImpressionNum            int    `json:"cur_impression_num"`
	ClickNum            int    `json:"click_num"`
	LikeNum            int    `json:"like_num"`
	CurLikeNum            int    `json:"cur_like_num"`
	CommentNum            int    `json:"comment_num"`
	CurCommentNum            int    `json:"cur_comment_num"`
	ShareNum            int    `json:"share_num"`
	FollowPhotoNum            int    `json:"follow_photo_num"`
	CurFollowPhotoNum            int    `json:"cur_follow_photo_num"`
}


//自增处理

func (m *PictureStatistics)IncData(column map[string]int) bool {
	if m.PId <= 0{
		return false
	}
	var updateGorm = make(map[string]interface{})
	tx := db.Model(&PictureStatistics{}).Where("p_id = ? ", m.PId)
	for k,v := range column{
		switch k {
		case "impression":
			updateGorm["impression_num"] = gorm.Expr("impression_num + ?",v)
			updateGorm["cur_impression_num"] = gorm.Expr("cur_impression_num + ?",v)
		case "click":
			updateGorm["click_num"] = gorm.Expr("click_num + ?",v)
		case "like":
			updateGorm["like_num"] = gorm.Expr("like_num + ?",v)
			updateGorm["cur_like_num"] = gorm.Expr("cur_like_num + ?",v)
		case "comment":
			updateGorm["comment_num"] = gorm.Expr("comment_num + ?",v)
			updateGorm["cur_comment_num"] = gorm.Expr("cur_comment_num + ?",v)
		case "share":
			updateGorm["share_num"] = gorm.Expr("share_num + ?",v)
		case "follow_photo":
			updateGorm["follow_photo_num"] = gorm.Expr("follow_photo_num + ?",v)
			updateGorm["cur_follow_photo_num"] = gorm.Expr("cur_follow_photo_num + ?",v)
		default:
		}
	}
	if len(updateGorm) > 0{
		if err := tx.Updates(updateGorm).Error;err != nil{
			global.Logger.Error("PictureStatisticsIncData"+err.Error())
			return false
		}
	}
	return true

}


//更新权重为初值值
func (m *PictureStatistics)InitWeightData(ids []int) error {
	if len(ids) > 0{
		if err := db.Model(&PictureStatistics{}).Where("p_id IN ? ", ids).Updates(PictureStatistics{CurImpressionNum: 1,CurLikeNum: 1,CurCommentNum: 1,CurFollowPhotoNum: 1}).Error; err != nil {
			global.Logger.Error("pictureStatisticsUpdateWeight"+err.Error())
			return err
		}

	}else{
		var err error
		return err
	}

	return nil
}

func (m *PictureStatistics)GetStatisticsKeyList(pIds []int) (map[int]PictureStatistics,error) {
	var pictureStatistics []PictureStatistics
	result := make(map[int]PictureStatistics)
	if len(pIds) > 0 {
		err := db.Where("p_id in (?)", pIds).Find(&pictureStatistics).Error
		if err != nil {
			fmt.Println("pictureStatistics,GetStatisticsKeyList err:", err)
		}
		for _, v := range pictureStatistics {
			result[v.PId] = v
		}
		return result,nil
	}
	return result,nil
}
