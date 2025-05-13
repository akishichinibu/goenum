package gen

func Gen(workdir string) (err error) {
	err = emitEnum(
		scanDecl(
			scanUnits(
				workdir,
			),
		),
	)

	return err
}
