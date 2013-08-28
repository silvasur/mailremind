package mysql

const (
	qUserByID = iota
	qUserByEmail
	qSetPWHash
	qSetActive
	qSetAcCode
	qDelUsersJobs
	qDelUser
	qGetOldInactiveUsers
	qCountJobs
	qJobsOfUser
	qJobFromUserAndID
	qSetSubject
	qSetContent
	qSetNext
	qDelJob
	qJobsBefore
	qInsertJob
	qInsertUser
	qSetChronos
	qEnd
)

const (
	qfragSelUser = "SELECT `id`, `email`, `passwd`, `location`, `active`, `activationcode`, `added` FROM `users` "
	qfragSelJob  = "SELECT `id`, `user`, `subject`, `content`, `next`, `chronos` FROM `jobs` "
)

var queries = map[int]string{
	qUserByID:            qfragSelUser + "WHERE `id` = ?",
	qUserByEmail:         qfragSelUser + "WHERE `email` = ?",
	qSetPWHash:           "UPDATE `users` SET `passwd` = ? WHERE `id` = ?",
	qSetActive:           "UPDATE `users` SET `active` = ? WHERE `id` = ?",
	qSetAcCode:           "UPDATE `users` SET `activationcode` = ? WHERE `id` = ?",
	qDelUsersJobs:        "DELETE FROM `jobs` WHERE `user` = ?",
	qDelUser:             "DELETE FROM `users` WHERE `id` = ?",
	qGetOldInactiveUsers: "SELECT `id` FROM `users` WHERE `active` = 0 AND `added` < ?",
	qCountJobs:           "SELECT COUNT(*) FROM `jobs` WHERE `user` = ?",
	qJobsOfUser:          qfragSelJob + "WHERE `user` = ?",
	qJobFromUserAndID:    qfragSelJob + "WHERE `user` = ? AND `id` = ?",
	qSetSubject:          "UPDATE `jobs` SET `subject` = ? WHERE `id` = ?",
	qSetContent:          "UPDATE `jobs` SET `content` = ? WHERE `id` = ?",
	qSetNext:             "UPDATE `jobs` SET `next` = ? WHERE `id` = ?",
	qDelJob:              "DELETE FROM `jobs` WHERE `id` = ?",
	qJobsBefore:          "SELECT `id` FROM `jobs` WHERE `next` <= ?",
	qInsertJob:           "INSERT INTO `jobs` (`user`, `subject`, `content`, `next`, `chronos`) VALUES (?, ?, ?, ?, ?)",
	qInsertUser:          "INSERT INTO `users` (`email`, `passwd`, `location`, `active`, `activationcode`, `added`) VALUES (?, ?, ?, ?, ?, ?)",
	qSetChronos:          "UPDATE `jobs` SET `chronos` = ? WHERE `id` = ?",
}
