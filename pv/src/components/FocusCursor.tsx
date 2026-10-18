import React from 'react';
import {useCurrentFrame} from 'remotion';
import {Cue, CursorLead, CursorRelease, CursorTrail, follow, Interactive} from '../lilt/spring';
import {C, FPS} from '../theme';

export type Rect = [x: number, y: number, w: number, h: number];
export type Move = {at: number; rect: Rect};
export type Press = {at: number; down: boolean};

type Props = {
	moves: Move[];
	presses?: Press[];
	// depth is how far each side squeezes in on a press.
	depth?: number;
	radius?: number;
};

// FocusCursor is ui.Cursor at video scale. Each edge has its own value: the
// edge facing the direction of travel eases with CursorLead and the far one
// with the slower CursorTrail, so the highlight stretches toward the target
// and contracts as it lands. It inverts what is under it, like the panel.
export const FocusCursor: React.FC<Props> = ({moves, presses = [], depth = 10, radius = 10}) => {
	const t = useCurrentFrame() / FPS;
	if (moves.length === 0 || t < moves[0].at) return null;

	const edges: Cue[][] = [[], [], [], []];
	let prev: Rect | null = null;
	for (const m of moves) {
		const [x, y, w, h] = m.rect;
		const target = [x, y, x + w, y + h];
		if (!prev) {
			// Grow out of the target's centre.
			const cx = x + w / 2;
			const cy = y + h / 2;
			[cx, cy, cx, cy].forEach((v, i) =>
				edges[i].push({at: m.at, jump: v, to: target[i], spring: CursorLead}),
			);
		} else {
			const dx = x + w / 2 - (prev[0] + prev[2] / 2);
			const dy = y + h / 2 - (prev[1] + prev[3] / 2);
			const lead = (d: number, hi: boolean) =>
				d === 0 ? CursorLead : d > 0 === hi ? CursorLead : CursorTrail;
			const springs = [lead(dx, false), lead(dy, false), lead(dx, true), lead(dy, true)];
			target.forEach((v, i) => edges[i].push({at: m.at, to: v, spring: springs[i]}));
		}
		prev = m.rect;
	}
	const [l, top, r, b] = edges.map((cues) => follow(t, cues));

	const press = follow(
		t,
		presses.map((p) => ({at: p.at, to: p.down ? 1 : 0, spring: p.down ? Interactive : CursorRelease})),
	);
	const inset = press * depth;
	const w = r - l - inset * 2;
	const h = b - top - inset;
	if (w < 1 || h < 1) return null;
	return (
		<div
			style={{
				position: 'absolute',
				left: l + inset,
				top: top + inset / 2,
				width: w,
				height: h,
				borderRadius: Math.min(radius, h / 2, w / 2),
				background: C.lit,
				mixBlendMode: 'difference',
				pointerEvents: 'none',
			}}
		/>
	);
};
