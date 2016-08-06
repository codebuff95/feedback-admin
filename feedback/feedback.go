package feedback

import(
)

type Feedback struct{
  Sectionid string  `bson:"sectionid"`
  Ratings []Rating  `bson:"ratings"`
  Addedon string  `bson:"addedon"`
}

type Rating struct{
  Facultyid string  `bson:"facultyid"`
  Subjectid string `bson:"subjectid"`
  Points []Point  `bson:"points"`
}

type Point struct{
  Questionid string `bson:"questionid"`
  Marks int `bson:"marks"`
}
