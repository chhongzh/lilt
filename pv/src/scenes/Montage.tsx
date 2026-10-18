import React from 'react';
import {AbsoluteFill, Sequence, useCurrentFrame} from 'remotion';
import {Panel} from '../components/Oled';
import {Code, Head} from '../components/Type';
import {Clicks, footageClicks} from '../components/Clicks';
import {mark} from './Footage';

export const CUT = 80;
const cuts = [
	{start: mark('sheet') - 6, head: 'Select 从底部升起 Sheet', code: 'ui.Sheet'},
	{start: mark('chips') + 4, head: '选中项下方的横条从中间长出', code: 'widget.Choice · widget.Toggles'},
	{start: mark('shelves'), head: '横向滚动，卡片跟着光标起伏', code: 'widget.ScrollHorizontal'},
	{start: mark('keyboard') - 4, head: '松手打开软键盘', code: 'widget.Input · widget.Keyboard'},
];
export const MONTAGE_LEN = CUT * cuts.length;

const PITCH = 10;

const Cut: React.FC<{start: number; head: string; code: string}> = ({start, head, code}) => {
	const f = useCurrentFrame();
	return (
		<AbsoluteFill>
			<Panel clip="reel" frame={start + f} pitch={PITCH} style={{position: 'absolute', left: (1920 - 128 * PITCH - 4 * PITCH * 2.2) / 2, top: 96}} />
			<div style={{position: 'absolute', left: 0, right: 0, top: 850, display: 'flex', flexDirection: 'column', alignItems: 'center'}}>
				<Head text={head} at={0.08} size={56} />
				<Code text={code} at={0.2} size={28} />
			</div>
			<Clicks clicks={footageClicks('reel', start, CUT)} />
		</AbsoluteFill>
	);
};

export const Montage: React.FC = () => (
	<AbsoluteFill>
		{cuts.map((c, i) => (
			<Sequence key={i} from={i * CUT} durationInFrames={i === cuts.length - 1 ? CUT * 2 : CUT}>
				<Cut {...c} />
			</Sequence>
		))}
	</AbsoluteFill>
);
