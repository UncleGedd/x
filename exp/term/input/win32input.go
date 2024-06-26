package input

import "github.com/erikgeiser/coninput"

func parseWin32InputKeyEvent(vkc coninput.VirtualKeyCode, _ coninput.VirtualKeyCode, r rune, keyDown bool, cks coninput.ControlKeyState, repeatCount uint16) Event {
	isCtrl := cks.Contains(coninput.LEFT_CTRL_PRESSED | coninput.RIGHT_CTRL_PRESSED)

	k, ok := vkKeyEvent[vkc]
	if !ok && isCtrl {
		k = vkCtrlRune(k, r, vkc)
	} else if !ok {
		k = KeyDownEvent{Rune: r}
	}
	if isCtrl {
		k.Mod |= Ctrl
	}
	if cks.Contains(coninput.LEFT_ALT_PRESSED | coninput.RIGHT_ALT_PRESSED) {
		k.Mod |= Alt
	}
	if cks.Contains(coninput.SHIFT_PRESSED) {
		k.Mod |= Shift
	}

	// XXX: the following keys when set mean that the key is ON, not that
	// it was pressed. We should probably ignore them.
	if cks.Contains(coninput.NUMLOCK_ON|coninput.CAPSLOCK_ON|coninput.SCROLLLOCK_ON) && k.Rune == 0 && k.Sym == 0 {
		return nil
	}

	var e Event = KeyDownEvent(k)
	k.IsRepeat = repeatCount > 1
	if !keyDown {
		e = KeyUpEvent(k)
	}

	if repeatCount <= 1 {
		return e
	}

	var kevents []Event
	for i := 0; i < int(repeatCount); i++ {
		kevents = append(kevents, e)
	}

	return MultiEvent(kevents)
}

var vkKeyEvent = map[coninput.VirtualKeyCode]KeyDownEvent{
	coninput.VK_RETURN:    {Sym: KeyEnter},
	coninput.VK_BACK:      {Sym: KeyBackspace},
	coninput.VK_TAB:       {Sym: KeyTab},
	coninput.VK_ESCAPE:    {Sym: KeyEscape},
	coninput.VK_SPACE:     {Sym: KeySpace, Rune: ' '},
	coninput.VK_UP:        {Sym: KeyUp},
	coninput.VK_DOWN:      {Sym: KeyDown},
	coninput.VK_RIGHT:     {Sym: KeyRight},
	coninput.VK_LEFT:      {Sym: KeyLeft},
	coninput.VK_HOME:      {Sym: KeyHome},
	coninput.VK_END:       {Sym: KeyEnd},
	coninput.VK_PRIOR:     {Sym: KeyPgUp},
	coninput.VK_NEXT:      {Sym: KeyPgDown},
	coninput.VK_DELETE:    {Sym: KeyDelete},
	coninput.VK_SELECT:    {Sym: KeySelect},
	coninput.VK_SNAPSHOT:  {Sym: KeyPrintScreen},
	coninput.VK_INSERT:    {Sym: KeyInsert},
	coninput.VK_LWIN:      {Sym: KeyLeftSuper},
	coninput.VK_RWIN:      {Sym: KeyRightSuper},
	coninput.VK_APPS:      {Sym: KeyMenu},
	coninput.VK_NUMPAD0:   {Sym: KeyKp0},
	coninput.VK_NUMPAD1:   {Sym: KeyKp1},
	coninput.VK_NUMPAD2:   {Sym: KeyKp2},
	coninput.VK_NUMPAD3:   {Sym: KeyKp3},
	coninput.VK_NUMPAD4:   {Sym: KeyKp4},
	coninput.VK_NUMPAD5:   {Sym: KeyKp5},
	coninput.VK_NUMPAD6:   {Sym: KeyKp6},
	coninput.VK_NUMPAD7:   {Sym: KeyKp7},
	coninput.VK_NUMPAD8:   {Sym: KeyKp8},
	coninput.VK_NUMPAD9:   {Sym: KeyKp9},
	coninput.VK_MULTIPLY:  {Sym: KeyKpMul},
	coninput.VK_ADD:       {Sym: KeyKpPlus},
	coninput.VK_SEPARATOR: {Sym: KeyKpComma},
	coninput.VK_SUBTRACT:  {Sym: KeyKpMinus},
	coninput.VK_DECIMAL:   {Sym: KeyKpPeriod},
	coninput.VK_DIVIDE:    {Sym: KeyKpDiv},
	coninput.VK_F1:        {Sym: KeyF1},
	coninput.VK_F2:        {Sym: KeyF2},
	coninput.VK_F3:        {Sym: KeyF3},
	coninput.VK_F4:        {Sym: KeyF4},
	coninput.VK_F5:        {Sym: KeyF5},
	coninput.VK_F6:        {Sym: KeyF6},
	coninput.VK_F7:        {Sym: KeyF7},
	coninput.VK_F8:        {Sym: KeyF8},
	coninput.VK_F9:        {Sym: KeyF9},
	coninput.VK_F10:       {Sym: KeyF10},
	coninput.VK_F11:       {Sym: KeyF11},
	coninput.VK_F12:       {Sym: KeyF12},
	coninput.VK_F13:       {Sym: KeyF13},
	coninput.VK_F14:       {Sym: KeyF14},
	coninput.VK_F15:       {Sym: KeyF15},
	coninput.VK_F16:       {Sym: KeyF16},
	coninput.VK_F17:       {Sym: KeyF17},
	coninput.VK_F18:       {Sym: KeyF18},
	coninput.VK_F19:       {Sym: KeyF19},
	coninput.VK_F20:       {Sym: KeyF20},
	coninput.VK_F21:       {Sym: KeyF21},
	coninput.VK_F22:       {Sym: KeyF22},
	coninput.VK_F23:       {Sym: KeyF23},
	coninput.VK_F24:       {Sym: KeyF24},
	coninput.VK_NUMLOCK:   {Sym: KeyNumLock},
	coninput.VK_SCROLL:    {Sym: KeyScrollLock},
	coninput.VK_LSHIFT:    {Sym: KeyLeftShift},
	coninput.VK_RSHIFT:    {Sym: KeyRightShift},
	coninput.VK_LCONTROL:  {Sym: KeyLeftCtrl},
	coninput.VK_RCONTROL:  {Sym: KeyRightCtrl},
	coninput.VK_LMENU:     {Sym: KeyLeftAlt},
	coninput.VK_RMENU:     {Sym: KeyRightAlt},
	coninput.VK_OEM_4:     {Rune: '['},
	// TODO: add more keys
}

func vkCtrlRune(k KeyDownEvent, r rune, kc coninput.VirtualKeyCode) KeyDownEvent {
	switch r {
	case '@':
		k.Rune = '@'
	case '\x01':
		k.Rune = 'a'
	case '\x02':
		k.Rune = 'b'
	case '\x03':
		k.Rune = 'c'
	case '\x04':
		k.Rune = 'd'
	case '\x05':
		k.Rune = 'e'
	case '\x06':
		k.Rune = 'f'
	case '\a':
		k.Rune = 'g'
	case '\b':
		k.Rune = 'h'
	case '\t':
		k.Rune = 'i'
	case '\n':
		k.Rune = 'j'
	case '\v':
		k.Rune = 'k'
	case '\f':
		k.Rune = 'l'
	case '\r':
		k.Rune = 'm'
	case '\x0e':
		k.Rune = 'n'
	case '\x0f':
		k.Rune = 'o'
	case '\x10':
		k.Rune = 'p'
	case '\x11':
		k.Rune = 'q'
	case '\x12':
		k.Rune = 'r'
	case '\x13':
		k.Rune = 's'
	case '\x14':
		k.Rune = 't'
	case '\x15':
		k.Rune = 'u'
	case '\x16':
		k.Rune = 'v'
	case '\x17':
		k.Rune = 'w'
	case '\x18':
		k.Rune = 'x'
	case '\x19':
		k.Rune = 'y'
	case '\x1a':
		k.Rune = 'z'
	case '\x1b':
		k.Rune = ']'
	case '\x1c':
		k.Rune = '\\'
	case '\x1f':
		k.Rune = '_'
	}

	switch kc {
	case coninput.VK_OEM_4:
		k.Rune = '['
	}

	return k
}
