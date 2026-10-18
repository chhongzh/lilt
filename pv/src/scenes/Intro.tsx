import React from 'react';
import {AbsoluteFill, useCurrentFrame} from 'remotion';
import {Panel} from '../components/Oled';
import {Keys} from '../components/Keys';
import {Body, Code, Head} from '../components/Type';
import {Clicks, footageClicks} from '../components/Clicks';

export const INTRO_LEN = 200;
const START = 14;

export const Intro: React.FC = () => {
	const f = useCurrentFrame();
	return (
		<AbsoluteFill>
			<Panel clip="reel" frame={START + f} pitch={8.6} style={{position: 'absolute', left: 120, top: 268}} />
			<div style={{position: 'absolute', left: 1330, top: 300, width: 520}}>
				<Head text="128 × 64" at={0.35} size={92} />
				<Body text="单色 OLED，一块 ESP32-S3。" at={0.55} />
				<Code text="examples/esp32s3-generic" at={0.75} size={26} />
				<Keys clip="reel" start={START} size={58} style={{marginTop: 64}} />
			</div>
			<Clicks clicks={footageClicks('reel', START, INTRO_LEN)} />
		</AbsoluteFill>
	);
};
