package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Article struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Title  string    `gorm:"not null;" json:"title"`
	Detail  string    `gorm:"not null;" json:"detail"`
	Thumbnail  string    `gorm:"size:255;not null" json:"thumbnail"`
	User    User      `json:"user"`
	UserID  uint32    `gorm:"not null" json:"user_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}


func (p *Article) Prepare() {
	p.ID = 0
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Detail = html.EscapeString(strings.TrimSpace(p.Detail))
	p.User = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Article) Validate() error {

	if p.Title == "" {
		return errors.New("Required Title")
	}
	if p.Detail == "" {
		return errors.New("Required Detail")
	}
	if p.Thumbnail == "" {
		return errors.New("Required Thumbnail")
	}
	if p.UserID < 1 {
		return errors.New("Required Author")
	}
	return nil
}

func (p *Article) SaveArticle(db *gorm.DB) (*Article, error) {
	var err error
	err = db.Debug().Model(&Article{}).Create(&p).Error
	if err != nil {
		return &Article{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&p.User).Error
		if err != nil {
			return &Article{}, err
		}
	}
	return p, nil
}

func (p *Article) FindAllArticles(db *gorm.DB) (*[]Article, error) {
	var err error
	posts := []Article{}
	err = db.Debug().Model(&Article{}).Limit(100).Find(&posts).Error
	if err != nil {
		return &[]Article{}, err
	}
	if len(posts) > 0 {
		for i, _ := range posts {
			err := db.Debug().Model(&User{}).Where("id = ?", posts[i].UserID).Take(&posts[i].User).Error
			if err != nil {
				return &[]Article{}, err
			}
		}
	}
	return &posts, nil
}

func (p *Article) FindArticleByID(db *gorm.DB, pid uint64) (*Article, error) {
	var err error
	err = db.Debug().Model(&Article{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Article{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&p.User).Error
		if err != nil {
			return &Article{}, err
		}
	}
	return p, nil
}

func (p *Article) UpdateAArticle(db *gorm.DB) (*Article, error) {

	var err error

	err = db.Debug().Model(&Article{}).Where("id = ?", p.ID).Updates(Article{Title: p.Title, Detail: p.Detail, Thumbnail: p.Thumbnail, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Article{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&p.User).Error
		if err != nil {
			return &Article{}, err
		}
	}
	return p, nil
}

func (p *Article) DeleteAArticle(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Article{}).Where("id = ? and author_id = ?", pid, uid).Take(&Article{}).Delete(&Article{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Post not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}