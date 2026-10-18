import React from 'react';
import {AbsoluteFill, useCurrentFrame} from 'remotion';
import {Panel} from '../components/Oled';
import {Keys} from '../components/Keys';
import {Body, Head} from '../components/Type';
import {Clicks, footageClicks} from '../components/Clicks';
import {index} from '../footage';
import {FPS} from '../theme';

export type Beat = {head: string; body: string; at: number};

type Props = {
	start: number;
	length: number;
	beats: Beat[];
	// side puts the panel on the left or the right, so consecutive scenes
	// do not share one composition.
	side: 'left' | 'right';
	// delay holds the first frame while the scene slides in, so the first
	// key lands once the screen has settled.
	delay?: number;
};

const PITCH = 8;
const PANEL_W = 128 * PITCH + 2 * PITCH * 2.2;

// Footage is a panel playing a stretch of the reel next to captions that
// change on cue. Each beat replaces the last one as the screen stack does:
// the old text rises out while the new one rises in.
export const Footage: React.FC<Props> = ({start, length, beats, side, delay = 0}) => {
	const f = Math.max(0, useCurrentFrame() - delay);
	const panelLeft = side === 'left' ? 110 : 1920 - 110 - PANEL_W;
	const textLeft = side === 'left' ? 110 + PANEL_W + 80 : 110;
	return (
		<AbsoluteFill>
			<Panel clip="reel" frame={start + f} pitch={PITCH} style={{position: 'absolute', left: panelLeft, top: 250}} />
			<div style={{position: 'absolute', left: textLeft, top: 250, width: 1920 - PANEL_W - 300}}>
				{beats.map((b, i) => {
					const out = beats[i + 1]?.at;
					return (
						<div key={i} style={{position: 'absolute', left: 0, top: 0, width: '100%'}}>
							<Head text={b.head} at={b.at} out={out} size={60} />
							<Body text={b.body} at={b.at + 0.12} out={out} size={36} />
						</div>
					);
				})}
				<Keys clip="reel" start={start - delay} size={52} style={{position: 'absolute', top: 330}} />
			</div>
			<Clicks clicks={footageClicks('reel', start, length - delay, delay)} />
		</AbsoluteFill>
	);
};

// at turns a reel frame into local seconds for a scene starting at start.
export const at = (start: number, reelFrame: number, delay = 0) => (reelFrame - start + delay) / FPS;
export const mark = (name: string) => index.reel.marks[name];
// keyAt finds the n-th event of a kind at or after a reel frame.
export const keyAt = (from: number, key: string, phase: string, n = 0) =>
	index.reel.keys.filter((e) => e.frame >= from && e.key === key && e.phase === phase)[n].frame;
