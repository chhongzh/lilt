import React from 'react';
import {AbsoluteFill, Sequence} from 'remotion';
import {Screen, Transition} from './components/Screen';
import {Opener} from './scenes/Opener';
import {Intro, INTRO_LEN} from './scenes/Intro';
import {Compare, COMPARE_LEN} from './scenes/Compare';
import {Springs, SPRINGS_LEN} from './scenes/Springs';
import {Slow, SLOW_LEN} from './scenes/Slow';
import {Action, BEAT_LEN, Press, Wave} from './scenes/Beats';
import {Montage, MONTAGE_LEN} from './scenes/Montage';
import {Facts, FACTS_H, FACTS_LEN} from './scenes/Facts';
import {Outro, OUTRO_LEN} from './scenes/Outro';
import {C} from './theme';

type Scene = {C: React.FC; len: number; enter: Transition; height?: number};

const scenes: Scene[] = [
	{C: Opener, len: 210, enter: 'none'},
	{C: Intro, len: INTRO_LEN, enter: 'slide'},
	{C: Compare, len: COMPARE_LEN, enter: 'slide'},
	{C: Springs, len: SPRINGS_LEN, enter: 'slide'},
	{C: Slow, len: SLOW_LEN, enter: 'slide'},
	{C: Press, len: BEAT_LEN, enter: 'slide'},
	{C: Action, len: BEAT_LEN, enter: 'slide'},
	{C: Wave, len: BEAT_LEN, enter: 'slide'},
	{C: Montage, len: MONTAGE_LEN, enter: 'slide'},
	{C: Facts, len: FACTS_LEN, enter: 'sheet', height: FACTS_H},
	{C: Outro, len: OUTRO_LEN, enter: 'slide'},
];

// A covered scene stays mounted until the next opaque scene has slid over
// it, the way the screen stack hides screens only once they are covered.
const TAIL = 60;

const starts = scenes.reduce<number[]>((acc, s, i) => [...acc, i === 0 ? 0 : acc[i - 1] + scenes[i - 1].len], []);
export const DURATION = starts[starts.length - 1] + scenes[scenes.length - 1].len;

const hiddenAt = (i: number) => {
	for (let j = i + 1; j < scenes.length; j++) if (scenes[j].enter !== 'sheet') return starts[j] + TAIL;
	return DURATION;
};

export const Video: React.FC = () => (
	<AbsoluteFill style={{background: C.bg}}>
		{scenes.map((s, i) => {
			const next = scenes[i + 1];
			return (
				<Sequence key={i} from={starts[i]} durationInFrames={hiddenAt(i) - starts[i]}>
					<Screen
						enter={starts[i]}
						cover={next ? starts[i + 1] : undefined}
						transition={s.enter}
						coveredBy={next?.enter}
						height={s.height}
					>
						<s.C />
					</Screen>
				</Sequence>
			);
		})}
	</AbsoluteFill>
);
