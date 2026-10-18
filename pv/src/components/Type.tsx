import React from 'react';
import {C, mono, sans} from '../theme';
import {WaveText} from './WaveText';

type Line = {text: string; at: number; out?: number};

// Lines splits on \n; each line starts its wave a beat after the one above.
const Lines: React.FC<Line & {stagger: number}> = ({text, at, out, stagger}) => (
	<>
		{text.split('\n').map((l, i) => (
			<WaveText key={i} text={l} at={at + i * 0.07} out={out === undefined ? undefined : out + i * 0.04} stagger={stagger} />
		))}
	</>
);

// Head is the scene's one focal line; Body sits under it in the dim tone.
export const Head: React.FC<Line & {size?: number}> = ({size = 76, ...l}) => (
	<div style={{fontFamily: sans, fontSize: size, fontWeight: 700, color: C.lit, lineHeight: 1.2}}>
		<Lines {...l} stagger={0.06} />
	</div>
);

export const Body: React.FC<Line & {size?: number}> = ({size = 40, ...l}) => (
	<div style={{fontFamily: sans, fontSize: size, fontWeight: 400, color: C.dim, lineHeight: 1.45, marginTop: size * 0.4}}>
		<Lines {...l} stagger={0.03} />
	</div>
);

// Code names a real identifier from the framework, in amber when it is a
// physics parameter.
export const Code: React.FC<Line & {size?: number; amber?: boolean; color?: string}> = ({size = 30, amber, color, ...l}) => (
	<div style={{fontFamily: mono, fontSize: size, color: color ?? (amber ? C.amber : C.dim), lineHeight: 1.5, marginTop: size * 0.4}}>
		<Lines {...l} stagger={0.02} />
	</div>
);
