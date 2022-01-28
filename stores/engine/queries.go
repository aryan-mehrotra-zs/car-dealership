package engine

const (
	insertEngine = "INSERT INTO engines (id,displacement,no_of_cylinder,`range`) VALUES (?,?,?,?)"
	getEngine    = "SELECT * FROM engines WHERE id=?"
	updateEngine = "UPDATE engines SET `displacement=?,no_of_cylinder=?,`range`=?` WHERE id=?"
	deleteEngine = "DELETE FROM engines WHERE id = ?;"
)
