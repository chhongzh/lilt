import React from 'react';
import {Audio, Sequence, staticFile} from 'remotion';
import {ClipName, index} from '../footage';

export type Click = {frame: number; down: boolean};

// Clicks plays a key sound per press and release, at local frames.
export const Clicks: React.FC<{clicks: Click[]}> = ({clicks}) => (
	<>
		{clicks.map((c, i) => (
			<Sequence key={i} from={c.frame} durationInFrames={8} layout="none">
				<Audio src={staticFile(c.down ? 'sfx/press.wav' : 'sfx/release.wav')} volume={c.down ? 0.5 : 0.35} />
			</Sequence>
		))}
	</>
);

// footageClicks lists the clip's key events that fall in the shown window.
export const footageClicks = (clip: ClipName, start: number, length: number, offset = 0): Click[] =>
	index[clip].keys
		.filter((e) => e.phase !== 'repeat' && e.frame >= start && e.frame < start + length)
		.map((e) => ({frame: e.frame - start + offset, down: e.phase === 'press'}));
