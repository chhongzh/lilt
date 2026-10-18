import React from 'react';
import {useCurrentFrame} from 'remotion';
import {CursorRelease, follow, Interactive} from '../lilt/spring';
import {ClipName, index} from '../footage';
import {C, FPS, mono} from '../theme';

type Props = {
	clip: ClipName;
	// start is the clip frame shown at local frame 0.
	start: number;
	size?: number;
	style?: React.CSSProperties;
};

const layout: [key: string, label: string, col: number, row: number][] = [
	['up', '↑', 1, 0],
	['left', '←', 0, 1],
	['ok', 'OK', 1, 1],
	['right', '→', 2, 1],
	['down', '↓', 1, 2],
	['back', 'BACK', 3.4, 2],
];

// Keys shows the real key events behind the footage, read from the
// capture index. A cap squeezes on press and wobbles back on release with
// the same springs the cursor uses, so you can see press and release apart.
export const Keys: React.FC<Props> = ({clip, start, size = 64, style}) => {
	const frame = useCurrentFrame();
	const now = start + frame;
	const t = frame / FPS;
	const events = index[clip].keys;
	const gap = size * 0.16;
	return (
		<div style={{position: 'relative', width: size * 4.4 + gap * 3, height: size * 3 + gap * 2, ...style}}>
			{layout.map(([key, label, col, row]) => {
				const cues = events
					.filter((e) => e.key === key && e.phase !== 'repeat' && e.frame <= now + 1)
					.map((e) => ({
						at: (e.frame - start) / FPS,
						to: e.phase === 'press' ? 1 : 0,
						spring: e.phase === 'press' ? Interactive : CursorRelease,
					}));
				const p = Math.max(0, follow(t, cues));
				const on = Math.min(1, p);
				const wide = label.length > 2;
				return (
					<div
						key={key}
						style={{
							position: 'absolute',
							left: col * (size + gap),
							top: row * (size + gap),
							width: wide ? size * 1.4 : size,
							height: size,
							borderRadius: size * 0.22,
							border: `2px solid ${on > 0.05 ? C.lit : C.edge}`,
							background: `rgba(200,230,255,${on * 0.9})`,
							color: on > 0.5 ? C.bg : C.dim,
							transform: `scale(${1 - p * 0.12})`,
							display: 'flex',
							alignItems: 'center',
							justifyContent: 'center',
							fontFamily: mono,
							fontSize: wide ? size * 0.28 : size * 0.4,
							fontWeight: 500,
							boxSizing: 'border-box',
						}}
					>
						{label}
					</div>
				);
			})}
		</div>
	);
};
