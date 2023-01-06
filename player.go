package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const moveSpeed float64 = 4
const moveSpeedSlow float64 = 2

// Player is the player character
type Player struct {
	*Entity

	// Bullets shot per second
	ShootSpeed    float64
	CanShoot      bool
	lastShootTime time.Time

	MoveHitbox   Collidable
	DamageHitbox Collidable
}

// NewPlayer creates a new player instance
func NewPlayer(position Vector) *Player {
	player := &Player{
		Entity: &Entity{
			Position: position,
		},
		ShootSpeed: 10,
		CanShoot:   true,
	}

	player.MoveHitbox = &RectangleHitbox{
		Size: Vector{X: 32, Y: 32},
		Hitbox: Hitbox{
			Position: Vector{X: -16, Y: -16},
			Owner:    player.Entity,
		},
	}

	player.DamageHitbox = &RectangleHitbox{
		Size: Vector{X: 16, Y: 16},
		Hitbox: Hitbox{
			Position: Vector{},
			Owner:    player.Entity,
		},
	}

	return player
}

// Start is called when the player is added to the game
func (player *Player) Start() {}

var gameFieldHitbox = &RectangleHitbox{
	Hitbox: Hitbox{
		Position: Vector{X: 32, Y: 32},
	},
	Size: PlayfieldSize.Minus(Vector{X: 64, Y: 64}),
}

// Update is called every game tick, and handles player behavior
func (player *Player) Update() {
	// Handle movement
	moveInput := Vector{
		X: AxisHorizontal.Get(0),
		Y: -AxisVertical.Get(0),
	}
	direction := moveInput.Angle()
	speed := 0.
	if moveInput.X != 0 || moveInput.Y != 0 {
		if ButtonSlow.Get(0) {
			speed = moveSpeedSlow
		} else {
			speed = moveSpeed
		}
	}

	// Allow sliding against walls
	for i := 0.; i < 60 && i > -60; i = -(i + Sign(i)) {
		mv := VectorFromAngle(direction + DegToRad(i)).ScaledBy(speed)
		if CollidesAt(player.MoveHitbox, player.Position.Plus(mv), gameFieldHitbox, Vector{}) {
			player.Velocity = mv
			player.Move(mv)
			break
		}
	}

	// Handle shooting
	if player.CanShoot && ButtonShoot.Get(0) {
		if time.Since(player.lastShootTime) > time.Second/time.Duration(player.ShootSpeed) {
			player.Shoot(
				player.Position.Copy().Minus(Vector{X: 0, Y: 0}),
				DegToRad(-90),
				6,
				25,
				0,
			)

			player.lastShootTime = time.Now()
		}
	}
}

// Die is called when the player dies
func (player *Player) Die() {
	// Make sure to clean up all the players bullets
	for obj := range BulletObjects {
		bullet, ok := obj.(*Bullet)
		if ok && bullet.Entity == *player.Entity {
			Destroy(bullet)
		}
	}
}

var (
	playerImage      = LoadImage("characters/player_forward.png", OriginCenter)
	playerLeftImage  = LoadImage("characters/player_left.png", OriginCenter)
	playerRightImage = LoadImage("characters/player_right.png", OriginCenter)
)

// Draw is called every frame to draw the player
func (player *Player) Draw(screen *ebiten.Image) {
	hAxis := AxisHorizontal.Get(0)
	image := playerImage

	if hAxis < 0 {
		image = playerLeftImage
	} else if hAxis > 0 {
		image = playerRightImage
	}

	image.Draw(screen, player.Position, Vector{X: 1, Y: 1}, 0)
}
