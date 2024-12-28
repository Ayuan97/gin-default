package models

import (
	"fmt"
	"gorm.io/gorm"
	"justus/global"
)

type PictureRecommend struct {
	Model
	PId    int     `json:"p_id" gorm:"column:p_id"`
	Lange  string  `json:"lange" gorm:"column:lange"`
	Weight float64 `json:"weight" gorm:"column:weight"`
}
type LangeRecommendCount struct {
	Lange string `json:"lange"`
	Count int    `json:"count"`
}

func (p *PictureRecommend) GetPictureRecommendListNot(NotPid []string, offset int, limit int, lange string) ([]*PictureRecommend, error) {
	var pictureRecommend []*PictureRecommend
	err := db.Not(map[string]interface{}{"p_id": NotPid}).Where("lange = ?", lange).Order("weight DESC").Offset(offset).Limit(limit).Find(&pictureRecommend).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return pictureRecommend, nil
}

//更新权重
func (p *PictureRecommend) UpdatePictureWeight() error {
	if err := db.Model(&PictureRecommend{}).Where("p_id = ? ", p.PId).Update("weight", p.Weight).Error; err != nil {
		global.Logger.Error("UpdatePictureWeight_recommend" + err.Error())
		return err
	}
	return nil
}

// 获取每个语言的推荐数量
func (p *PictureRecommend) GetPictureRecommendCount() ([]*LangeRecommendCount, error) {
	var langeRecommendCount []*LangeRecommendCount

	err := db.Model(&PictureRecommend{}).Select("count(id) as count,lange").Group("lange").Find(&langeRecommendCount).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return langeRecommendCount, nil
}

//查询过期的推荐信息 并删除
func (p *PictureRecommend) GetPictureRecommendExpire(offset int) bool {
	var pictureRecommend PictureRecommend
	err := db.Where("lange = ?", p.Lange).Order("created_at desc").Offset(offset).Limit(1).First(&pictureRecommend).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println(err)
	}
	//fmt.Println("超过的第一条信息id:",pictureRecommend.ID)
	//fmt.Println("p.Lange:",p.Lange)
	err = db.Where("id < ? AND lange = ?", pictureRecommend.ID, p.Lange).Delete(&PictureRecommend{}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println(err)
	}
	return true
}
