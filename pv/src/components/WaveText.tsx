import React from 'react';
import {useCurrentFrame} from 'remotion';
import {Default, follow, Spring} from '../lilt/spring';
import {FPS} from '../theme';

type Props = {
	text: string;
	// in and out are in local seconds; out is optional.
	at: number;
	out?: number;
	// stagger slows each glyph's spring a little more than the last, the
	// way widget.Scroll slows rows further from the focus.
	stagger?: number;
	spring?: Spring;
	style?: React.CSSProperties;
};

// WaveText rises into a clipped line glyph by glyph. Every glyph starts
// together; only the springs differ, so the line arrives as a wave.
export const WaveText: React.FC<Props> = ({text, at, out, stagger = 0.06, spring = Default, style}) => {
	const t = useCurrentFrame() / FPS;
	const glyphs = Array.from(text);
	return (
		<div style={{overflow: 'hidden', paddingBottom: '0.12em', marginBottom: '-0.12em', whiteSpace: 'pre', ...style}}>
			{glyphs.map((g, k) => {
				const s = {response: spring.response * (1 + stagger * k), damping: spring.damping};
				const cues = [{at, to: 0, spring: s}];
				if (out !== undefined) cues.push({at: out, to: -1.25, spring: s});
				const y = follow(t, cues, 1.25, s);
				return (
					<span key={k} style={{display: 'inline-block', transform: `translateY(${y}em)`}}>
						{g}
					</span>
				);
			})}
		</div>
	);
};
