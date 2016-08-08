# feedback-admin

The admin-side module for a feedback web application, which can be utilised at a college/university for collecting and managing feedbacks over a network, submitted by students for their respective teachers.

This application:

1. is made with Go (Google's Golang) as the primary server-side logic handling language.

2. uses 'github.com/codebuff95/uafm' as a user and form management library, and the uafm configuration files can be found in the directory '/feedbackadminres/uafmconfig.json'.

3. harnesses MongoDB's powers to help realise all the data-keeping tasks. This was made possible using the MongoDB driver for Go language, mgo (pronounced as mango).

4. has been designed giving data security a paramount importance. Maximum loopholes have been tried to be addressed to.
