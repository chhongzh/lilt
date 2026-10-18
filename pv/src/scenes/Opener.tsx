import React from 'react';
import {AbsoluteFill, useCurrentFrame} from 'remotion';
import logo from '../../public/footage/logo.json';
import {follow, wave} from '../lilt/spring';
import {FocusCursor} from '../components/FocusCursor';
import {Clicks} from '../components/Clicks';
import {WaveText} from '../components/WaveText';
import {C, FPS, H, sans} from '../theme';

const P = 40;
const LW = logo.W * P;
const LH = logo.H * P;
const LX = (1920 - LW) / 2;
const LY = 190;

// Wordmark draws "lilt" exactly as the panel font renders it, one dot per
// pixel. Columns rise on the list's wave: the centre lands first.
export const Wordmark: React.FC<{at: number; pitch: number; rise?: number}> = ({at, pitch, rise = H}) => {
	const t = useCurrentFrame() / FPS;
	const mid = (logo.W - 1) / 2;
	return (
		<div style={{position: 'relative', width: logo.W * pitch, height: logo.H * pitch, filter: `drop-shadow(0 0 ${pitch * 0.5}px rgba(200,230,255,0.55))`}}>
			{logo.Pixels.map(([x, y], i) => {
				const k = Math.round(Math.abs(x - mid));
				const dy = follow(t, [{at, to: 0, spring: wave(k, 0.25)}], rise);
				return (
					<div
						key={i}
						style={{
							position: 'absolute',
							left: x * pitch,
							top: y * pitch + dy,
							width: pitch * 0.86,
							height: pitch * 0.86,
							background: C.lit,
						}}
					/>
				);
			})}
		</div>
	);
};

const TAG_W = 820;
const TAG_Y = 700;

export const Opener: React.FC = () => {
	const pad = 44;
	return (
		<AbsoluteFill>
			<div style={{position: 'absolute', left: LX, top: LY}}>
				<Wordmark at={0.25} pitch={P} />
			</div>
			<div
				style={{
					position: 'absolute',
					left: (1920 - TAG_W) / 2,
					top: TAG_Y,
					width: TAG_W,
					textAlign: 'center',
					fontFamily: sans,
					fontSize: 58,
					fontWeight: 500,
					color: C.lit,
					display: 'flex',
					justifyContent: 'center',
				}}
			>
				<WaveText text="给 TinyGo 用的 GUI 框架" at={2.15} />
			</div>
			<FocusCursor
				radius={18}
				depth={14}
				moves={[
					{at: 1.35, rect: [LX - pad, LY - pad, LW + pad * 2 - P * 0.14, LH + pad * 2 - P * 0.14]},
					{at: 2.45, rect: [(1920 - TAG_W) / 2, TAG_Y - 18, TAG_W, 112]},
				]}
				presses={[
					{at: 1.8, down: true},
					{at: 1.98, down: false},
				]}
			/>
			<Clicks
				clicks={[
					{frame: Math.round(1.8 * FPS), down: true},
					{frame: Math.round(1.98 * FPS), down: false},
				]}
			/>
		</AbsoluteFill>
	);
};
