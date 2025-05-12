package gen

func Gen(workdir string) error {
	err :=
		emitEnum(
			scanDecl(
				scanUnits(
					workdir,
				),
			),
		)
	return err
}
