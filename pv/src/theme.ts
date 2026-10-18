import {loadFont} from '@remotion/google-fonts/NotoSansMono';

// The preview tool's OLED palette: the panel is the product, so its unlit
// black and lit ice-white are the whole page. Amber is the one accent and
// only ever marks physics: curves, parameters, the spring readouts.
export const C = {
	bg: '#080A0E',
	lit: '#C8E6FF',
	dim: '#8A9BB0',
	faint: '#141922',
	bezel: '#0C0F14',
	edge: '#1C222C',
	amber: '#FFB547',
};

// Noto Sans is what lilt itself draws with on the panel (notosans 12pt), so
// the captions are the same family. Noto Sans SC is installed system-wide.
export const sans = '"Noto Sans SC", "Noto Sans", "Microsoft YaHei", sans-serif';
export const mono = `${loadFont('normal', {weights: ['400', '500'], subsets: ['latin']}).fontFamily}, monospace`;

export const W = 1920;
export const H = 1080;
export const FPS = 60;
