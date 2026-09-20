package consts

import "testing"

func TestPodSaleModeConsts(t *testing.T) {
	t.Parallel()

	if PodAnnotationSaleModeKey != "katalyst.kubewharf.io/sale_mode" {
		t.Fatalf("unexpected sale mode annotation key: %q", PodAnnotationSaleModeKey)
	}
	if PodSaleModeSpot != "spot" {
		t.Fatalf("unexpected spot sale mode: %q", PodSaleModeSpot)
	}
	if PodSaleModeScheduled != "scheduled" {
		t.Fatalf("unexpected scheduled sale mode: %q", PodSaleModeScheduled)
	}
	if PodSaleModeReserved != "reserved" {
		t.Fatalf("unexpected reserved sale mode: %q", PodSaleModeReserved)
	}
	if PodSaleModeDefault != "default" {
		t.Fatalf("unexpected default sale mode: %q", PodSaleModeDefault)
	}
}
