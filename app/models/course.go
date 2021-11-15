package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Course struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"not null;" json:"title"`
	User      User      `json:"user"`
	UserID    uint32    `gorm:"not null" json:"user_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (c *Course) Prepare() {
	c.ID = 0
	c.Title = html.EscapeString(strings.TrimSpace(c.Title))
	c.User = User{}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

func (c *Course) Validate() error {
	if c.Title == "" {
		return errors.New("Required Title")
	}

	if c.UserID < 1 {
		return errors.New("Required Admin")
	}
	return nil
}

func (c *Course) SaveCourse(db *gorm.DB) (*Course, error) {
	var err error
	err = db.Debug().Model(&Course{}).Create(&c).Error
	if err != nil {
		return &Course{}, err
	}

	if c.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", c.UserID).Take(&c.User).Error
		if err != nil {
			return &Course{}, err
		}
	}

	return c, nil
}

func (c *Course) FindAllCourses(db *gorm.DB) (*[]Course, error) {
	var err error
	posts := []Course{}
	err = db.Debug().Model(&Course{}).Limit(100).Find(&posts).Error
	if err != nil {
		return &[]Course{}, err
	}

	if len(posts) > 0 {
		for i, _ := range posts {
			err := db.Debug().Model(&User{}).Where("id = ?", posts[i].UserID).Take(&posts[i].User).Error
			if err != nil {
				return &[]Course{}, err
			}
		}
	}

	return &posts, nil
}

func (c *Course) FindCourseByID(db *gorm.DB, cid uint64) (*Course, error) {
	var err error
	err = db.Debug().Model(&Course{}).Where("id = ?", cid).Take(&c).Error
	if err != nil {
		return &Course{}, err
	}

	if c.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", c.UserID).Take(&c.User).Error
		if err != nil {
			return &Course{}, err
		}
	}

	return c, nil
}

func (c *Course) UpdateCourse(db *gorm.DB) (*Course, error) {
	var err error
	err = db.Debug().Model(&Course{}).Where("id = ?", c.UserID).Updates(Course{Title: c.Title, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Course{}, err
	}

	if c.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", c.UserID).Take(&c.User).Error
		if err != nil {
			return &Course{}, err
		}
	}

	return c, nil
}

func (c *Course) DeleteCourse(db *gorm.DB, cid uint64, uid uint32) (int64, error) {
	db = db.Debug().Model(&Course{}).Where("id = ? and author_id = ?", cid, uid).Take(&Course{}).Delete(&Course{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("post not found")
		}
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
