#include <alsa/asoundlib.h>

// Works for a while but after a couple of seconds the connection to the
// PulseAudio server is lost. Kept for posteriority.
// 
// How to use in dwmstatus.go:
//
// // #cgo LDFLAGS: -lasound
// // #include "get_volume.h"
// import "C"
//
// func volume() string {
// 	v := C.get_volume()
// 	if v < 0 {
// 		return "Volume ERR"
// 	}
// 	return fmt.Sprintf("Volume %d%%", v)
// }
//
// Source: http://stackoverflow.com/questions/7657624/get-master-sound-volume-in-c-in-linux
int get_volume(void)
{
	static const char* mix_name = "Master";
	static const char* card = "default";
	static int mix_index = 0;

	snd_mixer_selem_id_t* sid;
	snd_mixer_selem_id_alloca(&sid);
	snd_mixer_selem_id_set_index(sid, mix_index);
	snd_mixer_selem_id_set_name(sid, mix_name);

	snd_mixer_t* handle;
	if ((snd_mixer_open(&handle, 0)) < 0) {
		return -1;
	}

	if ((snd_mixer_attach(handle, card)) < 0) {
		snd_mixer_close(handle);
		return -2;
	}

	if ((snd_mixer_selem_register(handle, NULL, NULL)) < 0) {
		snd_mixer_close(handle);
		return -3;
	}

	int ret = snd_mixer_load(handle);
	if (ret < 0) {
		snd_mixer_close(handle);
		return -4;
	}

	snd_mixer_elem_t* elem = snd_mixer_find_selem(handle, sid);
	if (!elem) {
		snd_mixer_close(handle);
		return -5;
	}

	long minv, maxv;
	snd_mixer_selem_get_playback_volume_range(elem, &minv, &maxv);

	long outvol;
	if (snd_mixer_selem_get_playback_volume(elem, 0, &outvol) < 0) {
		snd_mixer_close(handle);
		return -6;
	}

	outvol -= minv;
	maxv -= minv;
	return 100 * outvol / maxv;
}
