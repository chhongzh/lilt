import React from 'react';
import {Composition} from 'remotion';
import {DURATION, Video} from './Video';
import {FPS, H, W} from './theme';

export const Root: React.FC = () => (
	<Composition id="LiltPV" component={Video} durationInFrames={DURATION} fps={FPS} width={W} height={H} />
);
