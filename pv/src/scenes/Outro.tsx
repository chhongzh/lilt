import React from 'react';
import {AbsoluteFill, Img, staticFile, useCurrentFrame} from 'remotion';
import {follow} from '../lilt/spring';
import {FocusCursor} from '../components/FocusCursor';
import {Body, Code} from '../components/Type';
import {Clicks} from '../components/Clicks';
import {Wordmark} from './Opener';
import {C, FPS} from '../theme';

export const OUTRO_LEN = 330;

// The photo is the real board running the gallery; it drifts in on a very
// slow eased spring so even the last shot moves the way lilt moves.
export const Outro: React.FC = () => {
	const t = useCurrentFrame() / FPS;
	const drift = follow(t, [{at: 0, to: 1, spring: {response: 9, damping: 1, ease: true}}]);
	return (
		<AbsoluteFill>
			<div style={{position: 'absolute', left: 0, top: 0, width: 1440, height: 1080, overflow: 'hidden'}}>
				<Img
					src={staticFile('demo.jpg')}
					style={{
						width: 1440,
						height: 1080,
						objectFit: 'cover',
						transform: `scale(${1.12 - drift * 0.08}) translateX(${-drift * 30}px)`,
						transformOrigin: '40% 70%',
					}}
				/>
				<div style={{position: 'absolute', inset: 0, background: `linear-gradient(90deg, rgba(8,10,14,0) 45%, ${C.bg} 92%)`}} />
			</div>
			<div style={{position: 'absolute', left: 1240, top: 300}}>
				<Wordmark at={0.5} pitch={22} rise={700} />
				<div style={{marginTop: 60}}>
					<Body text="给 TinyGo 用的 GUI 框架" at={0.9} size={40} />
					<Code text="github.com/chhongzh/lilt" at={1.1} size={34} color={C.lit} />
					<Code text="MIT License" at={1.25} size={26} />
				</div>
			</div>
			<FocusCursor
				radius={12}
				depth={8}
				moves={[{at: 1.9, rect: [1222, 657, 580, 60]}]}
				presses={[
					{at: 2.9, down: true},
					{at: 3.08, down: false},
				]}
			/>
			<Clicks
				clicks={[
					{frame: Math.round(2.9 * FPS), down: true},
					{frame: Math.round(3.08 * FPS), down: false},
				]}
			/>
		</AbsoluteFill>
	);
};
