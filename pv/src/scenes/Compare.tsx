import React from 'react';
import {AbsoluteFill, useCurrentFrame} from 'remotion';
import {Panel} from '../components/Oled';
import {Keys} from '../components/Keys';
import {Body, Code, Head} from '../components/Type';
import {Clicks, footageClicks} from '../components/Clicks';
import {index} from '../footage';
import {C, FPS, sans} from '../theme';
import {WaveText} from '../components/WaveText';

export const COMPARE_LEN = 340;
// The clip starts at local DELAY, after the panels have arrived.
const DELAY = 30;
const PITCH = 5.8;
const SIDE = 150;

// Compare plays the same keys twice: once with anim.TimeScale at 40, so
// every spring lands in a frame or two, and once as lilt ships.
export const Compare: React.FC = () => {
	const f = useCurrentFrame();
	const clipFrame = f - DELAY;
	const end = index.compare.frames / FPS;
	const label = (text: string, at: number, color: string) => (
		<div style={{fontFamily: sans, fontSize: 40, fontWeight: 700, color, marginBottom: 22}}>
			<WaveText text={text} at={at} />
		</div>
	);
	return (
		<AbsoluteFill>
			<div style={{position: 'absolute', left: 0, right: 0, top: 70, display: 'flex', justifyContent: 'center'}}>
				<Head text="同一串按键，录了两遍" at={0.2} size={64} />
			</div>
			<div style={{position: 'absolute', left: SIDE, top: 230}}>
				{label('关掉弹簧', 0.35, C.dim)}
				<Panel clip="compare-stiff" frame={clipFrame} pitch={PITCH} />
				<Code text="anim.TimeScale = 40" at={0.5} size={26} />
			</div>
			<div style={{position: 'absolute', right: SIDE, top: 230}}>
				{label('lilt', 0.45, C.lit)}
				<Panel clip="compare" frame={clipFrame} pitch={PITCH} />
				<Code text="默认参数" at={0.6} size={26} />
			</div>
			<div style={{position: 'absolute', left: 0, right: 0, top: 790, display: 'flex', justifyContent: 'center'}}>
				<Keys clip="compare" start={-DELAY} size={40} />
			</div>
			<div style={{position: 'absolute', left: 0, right: 0, top: 960, display: 'flex', justifyContent: 'center'}}>
				<Body text="区别全在帧与帧之间。" at={DELAY / FPS + end - 0.5} size={44} />
			</div>
			<Clicks clicks={footageClicks('compare', 0, index.compare.frames, DELAY)} />
		</AbsoluteFill>
	);
};
