import React from 'react';
import {AbsoluteFill, useCurrentFrame} from 'remotion';
import {Bouncy, Default, follow, Smooth, Snappy, Spring, wave} from '../lilt/spring';
import {Body, Head} from '../components/Type';
import {C, FPS, mono, sans} from '../theme';

export const SPRINGS_LEN = 300;

// What each preset is for, from the comments in anim/spring.go.
const presets: [name: string, use: string, s: Spring][] = [
	['Default', '列表和布局', Default],
	['Smooth', '整屏转场', Smooth],
	['Snappy', '滑块和进度', Snappy],
	['Bouncy', '勾选标记弹入', Bouncy],
];

const FLIPS = [1.0, 2.05, 3.1, 4.15];
const TRACK_X = 640;
const TRACK_W = 820;
const DOT = 34;
const SPARK_X = 1560;
const SPARK_W = 220;
const SPARK_H = 96;
const SPARK_T = 0.9;

// curve samples the preset's step response, the same integration the
// framework runs, over SPARK_T seconds.
const curve = (s: Spring) => {
	const pts: string[] = [];
	for (let i = 0; i <= 90; i++) {
		const t = (i / 90) * SPARK_T;
		const v = follow(t, [{at: 0, to: 1, spring: s}]);
		pts.push(`${(i / 90) * SPARK_W},${SPARK_H - v * SPARK_H * 0.78}`);
	}
	return pts.join(' ');
};

export const Springs: React.FC = () => {
	const t = useCurrentFrame() / FPS;
	return (
		<AbsoluteFill>
			<div style={{position: 'absolute', left: 150, top: 96}}>
				<Head text="会动的东西，都挂在弹簧上" at={0.2} />
				<Body text="Response 是周期，Damping 是阻尼比，和 SwiftUI 一个写法。" at={0.4} />
			</div>
			{presets.map(([name, use, s], k) => {
				const y = 380 + k * 140;
				const rise = follow(t, [{at: 0.35, to: 0, spring: wave(k)}], 900);
				const cues = FLIPS.map((at, i) => ({at, to: i % 2 === 0 ? 1 : 0, spring: s}));
				const x = follow(t, cues);
				const drawn = Math.max(0, Math.min(1, (t - FLIPS[0]) / SPARK_T));
				return (
					<div key={name} style={{position: 'absolute', left: 0, top: y + rise, width: 1920, height: 110}}>
						<div style={{position: 'absolute', left: 150, top: 8, fontFamily: mono}}>
							<div style={{fontSize: 34, color: C.lit}}>anim.{name}</div>
							<div style={{fontSize: 26, color: C.amber, marginTop: 4}}>
								{s.response.toFixed(2)} / {s.damping.toFixed(2)}
							</div>
						</div>
						<div style={{position: 'absolute', left: 420, top: 18, fontFamily: sans, fontSize: 28, color: C.dim}}>{use}</div>
						<div
							style={{position: 'absolute', left: TRACK_X, top: 54, width: TRACK_W + DOT, height: 2, background: C.edge}}
						/>
						<div
							style={{
								position: 'absolute',
								left: TRACK_X + x * TRACK_W,
								top: 54 - DOT / 2 + 1,
								width: DOT,
								height: DOT,
								background: C.lit,
								boxShadow: '0 0 18px rgba(200,230,255,0.55)',
							}}
						/>
						<svg
							style={{position: 'absolute', left: SPARK_X, top: 4}}
							width={SPARK_W}
							height={SPARK_H + 4}
							viewBox={`0 -2 ${SPARK_W} ${SPARK_H + 4}`}
						>
							<line x1={0} x2={SPARK_W} y1={SPARK_H * 0.22} y2={SPARK_H * 0.22} stroke={C.edge} strokeWidth={2} strokeDasharray="4 6" />
							<polyline
								points={curve(s)}
								fill="none"
								stroke={C.amber}
								strokeWidth={3}
								pathLength={1}
								strokeDasharray={1}
								strokeDashoffset={1 - drawn}
							/>
						</svg>
					</div>
				);
			})}
			<div style={{position: 'absolute', left: 150, top: 960}}>
				<Body text="预设是指针：运行时改一处，用它的值一起变。" at={3.2} size={36} />
			</div>
		</AbsoluteFill>
	);
};
