/**
 * @Author: Sun
 * @Description:
 * @File:  model
 * @Version: 1.0.0
 * @Date: 2022/5/31 16:05
 */

package db

type CreateData struct {
	date  string `gorm:"date"`
	city  string `gorm:"city"`
	color string `gorm:"color"`
}
