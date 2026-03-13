package editor

// Media embedding support in PresentationEditor:
//   - Video: mp4, webm, avi, wmv, mov, mkv, m4v
//   - Audio: mp3, wav, wma, m4a, ogg, flac, aac
//
// Use AddVideo/AddVideoFromFile and AddAudio/AddAudioFromFile to place media
// on a slide. The editor writes both legacy relationship types
// (video/audio) and Office 2010+ p14:media extensions for compatibility.
// Video insertions accept custom poster frames and fall back to a built-in
// tiny PNG poster when none is provided.
// Audio insertions can optionally include a custom icon via
// AddAudioWithIcon/AddAudioWithIconFromFile.
//
// Playback option structs are available for parity API surface:
// NewVideoPlaybackOptions/NewAutoPlayVideoPlaybackOptions and
// NewAudioPlaybackOptions/NewAutoPlayAudioPlaybackOptions.
//
// Playback option APIs additionally emit slide timing media nodes so
// auto-play, loop, mute, volume, and across-slides audio settings are encoded
// in slide XML, including p14/p15 timing extension tags mapped to media rel IDs.
