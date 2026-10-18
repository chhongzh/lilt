import React from 'react';
import {AbsoluteFill, useCurrentFrame} from 'remotion';
import {follow, Smooth} from '../lilt/spring';
import {C, FPS, H, W} from '../theme';

export type Transition = 'slide' | 'sheet' | 'none';

type Props = {
	// enter and cover are global frames; cover is when the next scene arrives.
	enter: number;
	cover?: number;
	transition: Transition;
	coveredBy?: Transition;
	// height is the sheet's height; it rises to sit on the bottom edge.
	height?: number;
	children: React.ReactNode;
};

// Screen moves scenes the way ui.ScreenStack moves screens: SlideLeft
// pushes the new one in from the right on anim.Smooth while the old one
// backs off to a third of the width; Sheet rises from the bottom, sized to
// its content, over a screen that stays put.
export const Screen: React.FC<Props> = ({enter, cover, transition, coveredBy, height = H, children}) => {
	const g = useCurrentFrame() + enter;
	const t = g / FPS;
	const at = enter / FPS;
	const covered = cover !== undefined && coveredBy === 'slide';

	const xCues = [];
	if (transition === 'slide') xCues.push({at, jump: W, to: 0});
	if (covered) xCues.push({at: cover! / FPS, to: -W / 3});
	const x = follow(t, xCues, 0, Smooth);

	const y = transition === 'sheet' ? follow(t, [{at, jump: height + 2, to: 0}], 0, Smooth) : 0;
	const shade =
		cover !== undefined ? Math.max(0, Math.min(1, follow(t, [{at: cover / FPS, to: 0.6}], 0, Smooth))) : 0;

	const sheet = transition === 'sheet';
	return (
		<AbsoluteFill style={{transform: `translate(${x}px, ${y}px)`}}>
			<AbsoluteFill
				style={{
					top: sheet ? H - height : 0,
					height: sheet ? height : H,
					background: C.bg,
					borderTop: sheet ? `1px solid ${C.edge}` : undefined,
					borderRadius: sheet ? '28px 28px 0 0' : undefined,
					overflow: 'hidden',
				}}
			>
				{children}
			</AbsoluteFill>
			{shade > 0.001 && <AbsoluteFill style={{background: '#000', opacity: shade}} />}
		</AbsoluteFill>
	);
};
