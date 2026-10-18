import React, {useLayoutEffect, useRef} from 'react';
import {C} from '../theme';
import {ClipName, index, lit, useClip} from '../footage';

export type View = {x: number; y: number; w: number; h: number};
const FULL: View = {x: 0, y: 0, w: 128, h: 64};

type Props = {
	clip: ClipName;
	frame: number;
	// pitch is the on-video size of one panel pixel.
	pitch: number;
	view?: View;
	style?: React.CSSProperties;
};

// Oled draws one recorded frame as the panel shows it: square dots with a
// hairline gap, a faint matrix where pixels are off, and a soft bloom,
// since the panel emits light rather than reflecting it.
export const Oled: React.FC<Props> = ({clip, frame, pitch, view = FULL, style}) => {
	const data = useClip(clip);
	const ref = useRef<HTMLCanvasElement>(null);
	const f = Math.max(0, Math.min(index[clip].frames - 1, Math.round(frame)));
	const w = Math.round(view.w * pitch);
	const h = Math.round(view.h * pitch);

	useLayoutEffect(() => {
		const cv = ref.current;
		if (!cv || !data) return;
		const ctx = cv.getContext('2d')!;
		ctx.clearRect(0, 0, w, h);
		const dot = pitch * 0.86;
		const x0 = Math.floor(view.x);
		const y0 = Math.floor(view.y);
		const x1 = Math.min(128, Math.ceil(view.x + view.w));
		const y1 = Math.min(64, Math.ceil(view.y + view.h));
		const at = (x: number, y: number) => [(x - view.x) * pitch, (y - view.y) * pitch];

		if (pitch >= 6) {
			ctx.fillStyle = C.faint;
			for (let y = y0; y < y1; y++)
				for (let x = x0; x < x1; x++) {
					const [px, py] = at(x, y);
					ctx.fillRect(px, py, dot, dot);
				}
		}
		ctx.fillStyle = C.lit;
		for (let y = y0; y < y1; y++)
			for (let x = x0; x < x1; x++) {
				if (!lit(data, f, x, y)) continue;
				const [px, py] = at(x, y);
				ctx.fillRect(px, py, dot, dot);
			}
		ctx.save();
		ctx.globalCompositeOperation = 'lighter';
		ctx.globalAlpha = 0.45;
		ctx.filter = `blur(${Math.max(2, pitch * 0.9)}px)`;
		ctx.drawImage(cv, 0, 0);
		ctx.restore();
	}, [data, f, pitch, view.x, view.y, view.w, view.h, w, h]);

	return <canvas ref={ref} width={w} height={h} style={{display: 'block', ...style}} />;
};

// Panel is the module around the glass: a dark bezel with the active area
// inset, sized from the pitch so it scales with the footage.
export const Panel: React.FC<Props & {glass?: boolean}> = (p) => {
	const view = p.view ?? FULL;
	const pad = Math.max(10, p.pitch * 2.2);
	return (
		<div
			style={{
				padding: pad,
				background: C.bezel,
				border: `1px solid ${C.edge}`,
				borderRadius: Math.max(8, p.pitch * 1.2),
				width: view.w * p.pitch,
				boxSizing: 'content-box',
				...p.style,
			}}
		>
			<Oled {...p} style={undefined} />
		</div>
	);
};
