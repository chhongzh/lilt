// A port of lilt's anim package, so the video moves on the same physics as
// the framework: SwiftUI-style Response / Damping, semi-implicit Euler in
// 1/240 s substeps, and the same presets.

export type Spring = {response: number; damping: number; ease?: boolean};

export const Default: Spring = {response: 0.29, damping: 0.75};
export const Smooth: Spring = {response: 0.39, damping: 0.85};
export const Snappy: Spring = {response: 0.31, damping: 0.8};
export const Bouncy: Spring = {response: 0.25, damping: 0.45};
export const Interactive: Spring = {response: 0.14, damping: 1};

export const CursorLead: Spring = {response: 0.22, damping: 1, ease: true};
export const CursorTrail: Spring = {response: 0.52, damping: 1, ease: true};
export const CursorRelease: Spring = {response: 0.57, damping: 0.4};
export const CursorLeap: Spring = {response: 0.63, damping: 0.45};

// wave is the spring of the k-th row away from the focus in a staggered
// list, as widget.Scroll builds it.
export const wave = (k: number, stagger = 0.2, base: Spring = Default): Spring => ({
	response: base.response * (1 + stagger * Math.min(k, 7)),
	damping: base.damping,
});

const SUBSTEP = 1 / 240;
const rate = (s: Spring) => (2 * Math.PI) / s.response;

// A cue changes a value at a moment, in seconds: `to` retargets it (keeping
// its velocity, like Value.Set), `jump` teleports it (Value.Jump), `bump`
// throws it about that many units and lets it come back (Value.Bump).
export type Cue = {
	at: number;
	to?: number;
	jump?: number;
	bump?: number;
	spring?: Spring;
};

type State = {x: number; target: number; v: number; s: Spring};

const advance = (st: State, dt: number) => {
	if (dt <= 0) return;
	let x = st.x - st.target;
	let v = st.v;
	if (st.s.ease) {
		x *= Math.exp(-rate(st.s) * dt);
		v = 0;
	} else {
		const w = rate(st.s);
		const k = w * w;
		const c = 2 * st.s.damping * w;
		const n = Math.ceil(dt / SUBSTEP);
		const h = dt / n;
		for (let i = 0; i < n; i++) {
			v += (-k * x - c * v) * h;
			x += v * h;
		}
	}
	st.x = st.target + x;
	st.v = v;
};

// follow returns where a value starting at `initial` sits at time t after
// the cues, which must be sorted by time.
export const follow = (t: number, cues: Cue[], initial = 0, spring: Spring = Default): number => {
	const st: State = {x: initial, target: initial, v: 0, s: spring};
	let now = cues.length > 0 ? Math.min(cues[0].at, t) : t;
	for (const c of cues) {
		if (c.at > t) break;
		advance(st, c.at - now);
		now = c.at;
		if (c.spring) st.s = c.spring;
		if (c.jump !== undefined) {
			st.x = st.target = c.jump;
			st.v = 0;
		}
		if (c.to !== undefined) st.target = c.to;
		if (c.bump !== undefined && !st.s.ease) st.v += c.bump * rate(st.s) * Math.E;
	}
	advance(st, t - now);
	return st.x;
};

// to is the common case: rest at `from`, head for `to` at `at`.
export const to = (t: number, at: number, from: number, target: number, spring: Spring = Default) =>
	follow(t, [{at, to: target, spring}], from, spring);
