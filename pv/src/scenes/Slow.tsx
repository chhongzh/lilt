import React from 'react';
import {AbsoluteFill, useCurrentFrame} from 'remotion';
import {Panel} from '../components/Oled';
import {Keys} from '../components/Keys';
import {Code, Head} from '../components/Type';
import {Clicks, footageClicks} from '../components/Clicks';

export const SLOW_LEN = 330;
const START = 22;

// Slow is recorded with anim.TimeScale at 0.25, so the slow motion is the
// framework's own, not frames blended after the fact.
export const Slow: React.FC = () => {
	const f = useCurrentFrame();
	return (
		<AbsoluteFill>
			<Panel clip="slow" frame={START + f} pitch={10.5} style={{position: 'absolute', left: 150, top: 72}} />
			<div style={{position: 'absolute', left: 1600, top: 330}}>
				<Keys clip="slow" start={START} size={52} />
			</div>
			<div style={{position: 'absolute', left: 150, top: 840}}>
				<Head text="放慢 4 倍：光标先拉长，再收拢" at={0.3} size={64} />
				<Code text="CursorLead 0.22s · CursorTrail 0.52s · anim.TimeScale = 0.25" at={0.55} amber />
			</div>
			<Clicks clicks={footageClicks('slow', START, SLOW_LEN)} />
		</AbsoluteFill>
	);
};
