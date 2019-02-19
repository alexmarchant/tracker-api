package main

import (
  "time"
)

type Day struct {
  ID uint `gorm:"primary_key" json:"id"`
  CreatedAt time.Time `json:"createdAt"`
  UpdatedAt time.Time `json:"updatedAt"`
  DeletedAt *time.Time `json:"-"`
  Date string `gorm:"UNIQUE_INDEX;UNIQUE;NOT NULL" json:"date"`
  Bmr *int `json:"bmr,omitempty"`
  CaloriesIn *int `json:"caloriesIn,omitempty"`
  CaloriesOut *int `json:"caloriesOut,omitempty"`
  CaloriesGoal *int `json:"caloriesGoal,omitempty"`
  MilesRun *int `json:"milesRun,omitempty"`
  MilesRunGoal *int `json:"milesRunGoal,omitempty"`
  Drinks *int `json:"drinks,omitempty"`
  DrinksGoal *int `json:"drinksGoal,omitempty"`
}
