package command

type (
	VideoBinaryInsertFn     func(int, []byte, []byte, string, float64, float64, float64, float64) (int, error)
	VideoPathInsertFn       func(int, string, string, string, float64, float64, float64, float64) (int, error)
	AudioBinaryInsertFn     func(int, []byte, string, float64, float64, float64, float64) (int, error)
	AudioBinaryIconInsertFn func(int, []byte, []byte, string, float64, float64, float64, float64) (int, error)
	AudioPathInsertFn       func(int, string, string, float64, float64, float64, float64) (int, error)
	AudioPathIconInsertFn   func(int, string, string, string, float64, float64, float64, float64) (int, error)
	OLEBinaryInsertFn       func(int, []byte, []byte, string, float64, float64, float64, float64) (int, error)
	OLEPathInsertFn         func(int, string, string, string, float64, float64, float64, float64) (int, error)
)

func AdaptVideoBinaryInsert(insert VideoBinaryInsertFn) func(MediaPlacement, string, []byte, []byte) (int, error) {
	return func(placement MediaPlacement, mimeType string, videoData []byte, posterData []byte) (int, error) {
		return insert(
			placement.SlideIndex,
			videoData,
			posterData,
			mimeType,
			placement.X,
			placement.Y,
			placement.W,
			placement.H,
		)
	}
}

func AdaptVideoPathInsert(insert VideoPathInsertFn) func(MediaPlacement, string, string, string) (int, error) {
	return func(placement MediaPlacement, mimeType string, videoPath string, posterPath string) (int, error) {
		return insert(
			placement.SlideIndex,
			videoPath,
			posterPath,
			mimeType,
			placement.X,
			placement.Y,
			placement.W,
			placement.H,
		)
	}
}

func AdaptAudioBinaryInsert(insert AudioBinaryInsertFn) func(MediaPlacement, string, []byte, []byte) (int, error) {
	return func(placement MediaPlacement, mimeType string, audioData []byte, _ []byte) (int, error) {
		return insert(
			placement.SlideIndex,
			audioData,
			mimeType,
			placement.X,
			placement.Y,
			placement.W,
			placement.H,
		)
	}
}

func AdaptAudioPathInsert(insert AudioPathInsertFn) func(MediaPlacement, string, string, string) (int, error) {
	return func(placement MediaPlacement, mimeType string, audioPath string, _ string) (int, error) {
		return insert(
			placement.SlideIndex,
			audioPath,
			mimeType,
			placement.X,
			placement.Y,
			placement.W,
			placement.H,
		)
	}
}

func AdaptAudioBinaryInsertWithOptionalIcon(
	insert AudioBinaryInsertFn,
	insertWithIcon AudioBinaryIconInsertFn,
) func(MediaPlacement, string, []byte, []byte) (int, error) {
	return func(placement MediaPlacement, mimeType string, audioData []byte, iconData []byte) (int, error) {
		if len(iconData) > 0 {
			return insertWithIcon(
				placement.SlideIndex,
				audioData,
				iconData,
				mimeType,
				placement.X,
				placement.Y,
				placement.W,
				placement.H,
			)
		}
		return insert(
			placement.SlideIndex,
			audioData,
			mimeType,
			placement.X,
			placement.Y,
			placement.W,
			placement.H,
		)
	}
}

func AdaptAudioPathInsertWithOptionalIcon(
	insert AudioPathInsertFn,
	insertWithIcon AudioPathIconInsertFn,
) func(MediaPlacement, string, string, string) (int, error) {
	return func(placement MediaPlacement, mimeType string, audioPath string, iconPath string) (int, error) {
		if iconPath != "" {
			return insertWithIcon(
				placement.SlideIndex,
				audioPath,
				iconPath,
				mimeType,
				placement.X,
				placement.Y,
				placement.W,
				placement.H,
			)
		}
		return insert(
			placement.SlideIndex,
			audioPath,
			mimeType,
			placement.X,
			placement.Y,
			placement.W,
			placement.H,
		)
	}
}

func AdaptOLEBinaryInsert(insert OLEBinaryInsertFn) func(MediaPlacement, string, []byte, []byte) (int, error) {
	return func(placement MediaPlacement, progID string, objectData []byte, iconData []byte) (int, error) {
		return insert(
			placement.SlideIndex,
			objectData,
			iconData,
			progID,
			placement.X,
			placement.Y,
			placement.W,
			placement.H,
		)
	}
}

func AdaptOLEPathInsert(insert OLEPathInsertFn) func(MediaPlacement, string, string, string) (int, error) {
	return func(placement MediaPlacement, progID string, objectPath string, iconPath string) (int, error) {
		return insert(
			placement.SlideIndex,
			objectPath,
			iconPath,
			progID,
			placement.X,
			placement.Y,
			placement.W,
			placement.H,
		)
	}
}

func NewVideoInsertSpec(
	maxLen int,
	insertBinary func(MediaPlacement, string, []byte, []byte) (int, error),
	insertPath func(MediaPlacement, string, string, string) (int, error),
) MediaInsertSpec {
	return MediaInsertSpec{
		MetaKey:          KeyMimeType,
		PrimaryPathKey:   KeyPath,
		PrimaryDataKey:   KeyData,
		SecondaryPathKey: "poster_path",
		SecondaryDataKey: "poster_data",
		PrimaryMaxLen:    maxLen,
		SecondaryMaxLen:  maxLen,
		PrimaryLabel:     labelVideo,
		SecondaryLabel:   labelPoster,
		InsertBinary:     insertBinary,
		InsertPath:       insertPath,
	}
}

func NewAudioInsertSpec(
	maxLen int,
	insertBinary func(MediaPlacement, string, []byte, []byte) (int, error),
	insertPath func(MediaPlacement, string, string, string) (int, error),
) MediaInsertSpec {
	return MediaInsertSpec{
		MetaKey:          KeyMimeType,
		PrimaryPathKey:   KeyPath,
		PrimaryDataKey:   KeyData,
		SecondaryPathKey: "icon_path",
		SecondaryDataKey: "icon_data",
		PrimaryMaxLen:    maxLen,
		SecondaryMaxLen:  maxLen,
		PrimaryLabel:     "audio",
		SecondaryLabel:   labelIcon,
		InsertBinary:     insertBinary,
		InsertPath:       insertPath,
	}
}

func NewOLEInsertSpec(
	maxLen int,
	insertBinary func(MediaPlacement, string, []byte, []byte) (int, error),
	insertPath func(MediaPlacement, string, string, string) (int, error),
) MediaInsertSpec {
	return MediaInsertSpec{
		MetaKey:          "prog_id",
		PrimaryPathKey:   KeyPath,
		PrimaryDataKey:   KeyData,
		SecondaryPathKey: "icon_path",
		SecondaryDataKey: "icon_data",
		PrimaryMaxLen:    maxLen,
		SecondaryMaxLen:  maxLen,
		PrimaryLabel:     "object",
		SecondaryLabel:   labelIcon,
		InsertBinary:     insertBinary,
		InsertPath:       insertPath,
	}
}
