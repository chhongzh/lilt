import React from 'react';
import {at, Footage, keyAt, mark} from './Footage';

export const BEAT_LEN = 280;

const PRESS = 300;
export const Press: React.FC = () => (
	<Footage
		start={PRESS}
		length={BEAT_LEN}
		side="left"
		delay={30}
		beats={[
			{head: '按下只给反馈\n松手才触发', body: '按下 OK 光标收缩，\n松开带着回弹恢复。', at: 0.35},
			{head: '按住确认', body: '光标从左往右缩短，\n缩完才执行。', at: at(PRESS, keyAt(mark('hold'), 'ok', 'press'), 30) - 0.1},
		]}
	/>
);

const ACTION = 560;
export const Action: React.FC = () => (
	<Footage
		start={ACTION}
		length={BEAT_LEN}
		side="right"
		delay={20}
		beats={[
			{head: '加载，然后打勾', body: 'OnAction 在后台跑，\n按钮自己切换状态。', at: 0.35},
			{head: '失败就打叉\n再抖一下', body: 'OnAction 返回 error，\n状态机不用自己写。', at: at(ACTION, keyAt(mark('action'), 'ok', 'press', 1), 20) + 0.4},
		]}
	/>
);

const WAVE = mark('wave') - 40;
export const Wave: React.FC = () => (
	<Footage
		start={WAVE}
		length={BEAT_LEN}
		side="left"
		beats={[
			{head: '整列像波浪', body: '离光标越远的行，\n跟得越慢。', at: 0.35},
			{head: '到底了还按', body: '整列被推一下，\n松手弹回。', at: at(WAVE, keyAt(mark('wave'), 'down', 'release')) - 0.7},
		]}
	/>
);
