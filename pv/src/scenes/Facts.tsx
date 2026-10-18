import React from 'react';
import {AbsoluteFill} from 'remotion';
import {FocusCursor} from '../components/FocusCursor';
import {Clicks} from '../components/Clicks';
import {Body} from '../components/Type';
import {WaveText} from '../components/WaveText';
import {C, FPS, sans} from '../theme';

export const FACTS_LEN = 250;
export const FACTS_H = 640;

// From the README's rendering notes; each one is how the code is written.
const facts = ['每帧零分配', '只用 float32', '不依赖 strings 包'];
const ROW = 128;
const TOP = 170;
const LEFT = 150;

// Facts is a sheet over the montage; the cursor steps down the rows as if
// someone were pressing Down on the panel.
export const Facts: React.FC = () => {
	const steps = [0.7, 1.35, 2.0];
	return (
		<AbsoluteFill>
			<div style={{position: 'absolute', left: LEFT, top: 64}}>
				<Body text="写给单片机的约束" at={0.3} size={36} />
			</div>
			{facts.map((s, i) => (
				<div
					key={s}
					style={{position: 'absolute', left: LEFT + 40, top: TOP + i * ROW, fontFamily: sans, fontSize: 84, fontWeight: 700, color: C.lit, lineHeight: `${ROW - 18}px`}}
				>
					<WaveText text={s} at={0.35 + i * 0.08} />
				</div>
			))}
			<FocusCursor
				radius={20}
				moves={steps.map((at, i) => ({at, rect: [LEFT, TOP + i * ROW - 6, 1000, ROW - 6]}))}
				presses={steps.slice(1).flatMap((at) => [
					{at: at - 0.16, down: true},
					{at: at, down: false},
				])}
				depth={8}
			/>
			<Clicks clicks={steps.slice(1).flatMap((at) => [{frame: Math.round((at - 0.16) * FPS), down: true}, {frame: Math.round(at * FPS), down: false}])} />
		</AbsoluteFill>
	);
};
