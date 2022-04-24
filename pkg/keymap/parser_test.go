package keymap

import (
	"strings"
	"testing"
)

var cradioKeymapFile = `
#include <behaviors.dtsi>
#include <dt-bindings/zmk/keys.h>
#include <dt-bindings/zmk/bt.h>


&mt {
    flavor = "tap-preferred";
    tapping_term_ms = <200>;
};

/ {

    combos {
        compatible = "zmk,combos";
        combo_esc {
            timeout-ms = <50>;
            key-positions = <0 1>;
            bindings = <&kp ESC>;
        };

         combo_tab {
            timeout-ms = <50>;
            key-positions = <10 11>;
            bindings = <&kp TAB>;
        };

            combo_ralt {
            timeout-ms = <50>;
            key-positions = <17 16>;
            bindings = <&kp RALT>;
        };

                    combo_lalt {
            timeout-ms = <50>;
            key-positions = <11 12>;
            bindings = <&kp LALT>;
        };

                           combo_lgui {
            timeout-ms = <50>;
            key-positions = <12 13>;
            bindings = <&kp LGUI>;
        };


           combo_rgui {
            timeout-ms = <50>;
            key-positions = <17 18>;
            bindings = <&kp RGUI>;
        };



    };

        keymap {
                compatible = "zmk,keymap";

       		default_layer {
		bindings = <
		&kp Q &kp W &kp E &kp R &kp T 						&kp Y &kp U  &kp I    &kp O   &kp P
		&kp A &kp S &kp D &kp F &kp G 			        	        &kp H &kp J &kp K &kp L &kp QUOT
		&mt LSFT Z &kp X &kp C &kp V &kp B				        &kp N &kp M  &kp CMMA &kp DOT &mt LSFT RET
          				&mo 1 &kp LCTL  				&kp SPC &mo 2
		>;
		};

       		left_layer {
		bindings = <
		&kp NUM_1  &kp NUM_2    &kp NUM_3    &kp NUM_4    &kp NUM_5		&kp NUM_6 &kp NUM_7 &kp NUM_8 &kp NUM_9 &kp NUM_0
		&kp LC(LS(TAB))    &kp LC(S)    &kp DQT      &kp PIPE2     &kp HASH 		&kp MINUS &kp EQL  &kp LBKT &kp RBKT  &kp DEL
		&kp ESC    &kp TILDE 	&kp NON_US_BSLH &kp NON_US_HASH  &kp TILDE2 	&kp MINUS &kp GRAVE &kp LBKT &kp RBKT  &kp DEL
					    &mo 1  &kp LGUI  					&kp RGUI &mo 2
		>;
		};

		right_layer {
		bindings = <
		&kp BANG  &kp ATSN &kp HASH  &kp DLLR  &kp PRCT    			&kp CRRT  &kp AMPS &kp KMLT &kp LPRN &kp RPRN
		&kp HASH  &kp QMARK  &kp FSLH  &kp COLN  &kp SCLN 			&kp MINUS &kp KP_EQUAL  &kp LBRC  &kp RBRC   &kp BKSP
		&kp LSFT  &kp KPLS &kp LBKT &kp RBKT   &kp BSLH    			&kp UNDER &kp LEFT &kp DOWN &kp UP  &kp RIGHT
					  &mo 3 &kp LCTL  				&kp SPC  &mo 2
		>;
		};

		tri_layer {
		bindings = <
		&kp NUM_1  &kp NUM_2    &kp NUM_3    &kp NUM_4    &kp NUM_5 			&trans &trans   &trans   &trans  &trans
		&kp F1 &kp F2 &kp F3 &kp F4 &kp F5  						&trans &kp PG_UP  &kp K_VOL_UP &kp K_MUTE &trans
		&bt BT_CLR  &bt BT_NXT &bt BT_PRV &kp F6 &kp F7   				&trans &kp PG_DN  &kp K_VOL_DN  &trans &trans
					 &trans &trans  					&trans &trans
		>;
		};

	};
};
`

var sofleKeymapFile = `/*
* Copyright (c) 2020 The ZMK Contributors
*
* SPDX-License-Identifier: MIT
*/

#include <behaviors.dtsi>
#include <dt-bindings/zmk/keys.h>
#include <dt-bindings/zmk/bt.h>

/ {
	 keymap {
			 compatible = "zmk,keymap";

			 default_layer {
					 bindings = <
&kp GRAVE &kp N1 &kp N2   &kp N3   &kp N4    &kp N5                      &kp N6 &kp N7    &kp N8    &kp N9   &kp N0   &none
&kp ESC   &kp Q  &kp W    &kp E    &kp R     &kp T                       &kp Y  &kp U     &kp I     &kp O    &kp P    &kp BSPC
&kp TAB   &kp A  &kp S    &kp D    &kp F     &kp G                       &kp H  &kp J     &kp K     &kp L    &kp SEMI &kp SQT
&kp LSHFT &kp Z  &kp X    &kp C    &kp V     &kp B  &kp C_MUTE &none     &kp N  &kp M     &kp COMMA &kp DOT  &kp FSLH &kp RSHFT
								&kp LGUI &kp LALT &kp LCTRL &mo 1  &kp RET    &kp SPACE &mo 2  &kp RCTRL &kp RALT  &kp RGUI
					 >;

					 sensor-bindings = <&inc_dec_kp C_VOL_UP C_VOL_DN &inc_dec_kp PG_UP PG_DN>;
			 };

			 lower_layer {
					 bindings = <
&trans    &kp F1    &kp F2    &kp F3      &kp F4    &kp F5                    &kp F6    &kp F7   &kp F8          &kp F9    &kp F10   &kp F11
&kp GRAVE &kp N1    &kp N2    &kp N3      &kp N4    &kp N5                    &kp N6    &kp N7   &kp N8          &kp N9    &kp N0    &kp F12
&trans    &kp EXCL  &kp AT    &kp HASH    &kp DLLR  &kp PRCNT                 &kp CARET &kp AMPS &kp KP_MULTIPLY &kp LPAR  &kp RPAR  &kp PIPE
&trans    &kp EQUAL &kp MINUS &kp KP_PLUS &kp LBRC  &kp RBRC  &trans   &trans &kp LBKT  &kp RBKT &kp SEMI        &kp COLON &kp BSLH  &trans
									 &trans    &trans      &trans    &trans    &trans   &trans &trans    &trans   &trans          &trans
					 >;

					 sensor-bindings = <&inc_dec_kp C_VOL_UP C_VOL_DN &inc_dec_kp PG_UP PG_DN>;
			 };

			 raise_layer {
					 bindings = <
&bt BT_CLR &bt BT_SEL 0 &bt BT_SEL 1 &bt BT_SEL 2 &bt BT_SEL 3 &bt BT_SEL 4             &trans    &trans    &trans   &trans    &trans  &trans
&trans     &kp INS      &kp PSCRN    &kp K_CMENU  &trans       &trans                   &kp PG_UP &trans    &kp UP   &trans    &kp N0  &trans
&trans     &kp LALT     &kp LCTRL    &kp LSHFT    &trans       &kp CLCK                 &kp PG_DN &kp LEFT  &kp DOWN &kp RIGHT &kp DEL &kp BSPC
&trans     &kp K_UNDO   &kp K_CUT    &kp K_COPY   &kp K_PASTE  &trans  &trans   &trans  &trans    &trans    &trans   &trans    &trans  &trans
											 &trans       &trans       &trans       &trans  &trans   &trans  &trans    &trans    &trans   &trans
					 >;

					 sensor-bindings = <&inc_dec_kp C_VOL_UP C_VOL_DN &inc_dec_kp PG_UP PG_DN>;
			 };
	 };
};
`

func Test(t *testing.T) {
	cradioRaw := strings.NewReader(cradioKeymapFile)
	_, cradioErr := Parse(cradioRaw)
	if cradioErr != nil {
		t.Fatal(cradioErr)
	}

	sofleRaw := strings.NewReader(sofleKeymapFile)
	_, sofleErr := Parse(sofleRaw)
	if sofleErr != nil {
		t.Fatal(sofleErr)
	}
}
