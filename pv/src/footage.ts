import {useEffect, useState} from 'react';
import {cancelRender, continueRender, delayRender, staticFile} from 'remotion';
import reel from '../public/footage/reel.json';
import compare from '../public/footage/compare.json';
import compareStiff from '../public/footage/compare-stiff.json';
import slow from '../public/footage/slow.json';

// Footage is recorded by `go run ./capture`: one packed 1-bit bitmap per
// frame at 60 fps, plus an index of marks and key events.
export type ClipName = 'reel' | 'compare' | 'compare-stiff' | 'slow';

export type KeyEvent = {frame: number; key: string; phase: 'press' | 'repeat' | 'release'};
export type Index = {
	fps: number;
	width: number;
	height: number;
	frames: number;
	marks: Record<string, number>;
	keys: KeyEvent[];
};

export const index: Record<ClipName, Index> = {
	reel: reel as Index,
	compare: compare as Index,
	'compare-stiff': compareStiff as Index,
	slow: slow as Index,
};

export const FRAME_BYTES = 128 * 64 / 8;

const loads = new Map<ClipName, Promise<Uint8Array>>();
const loaded = new Map<ClipName, Uint8Array>();

const load = (name: ClipName) => {
	let p = loads.get(name);
	if (!p) {
		p = fetch(staticFile(`footage/${name}.bin`))
			.then((r) => r.arrayBuffer())
			.then((b) => {
				const data = new Uint8Array(b);
				loaded.set(name, data);
				return data;
			});
		loads.set(name, p);
	}
	return p;
};

export const useClip = (name: ClipName): Uint8Array | null => {
	const [data, setData] = useState(() => loaded.get(name) ?? null);
	const [handle] = useState(() => (loaded.has(name) ? null : delayRender(`footage ${name}`)));
	useEffect(() => {
		if (data) return;
		load(name)
			.then((d) => {
				setData(d);
				if (handle !== null) continueRender(handle);
			})
			.catch((e) => cancelRender(e));
	}, [name, data, handle]);
	return data;
};

export const lit = (data: Uint8Array, frame: number, x: number, y: number) => {
	const i = y * 128 + x;
	return (data[frame * FRAME_BYTES + (i >> 3)] & (0x80 >> (i & 7))) !== 0;
};
