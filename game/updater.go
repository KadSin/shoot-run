package game

import (
	"kadsin/shoot-run/game/assets"
	"kadsin/shoot-run/game/entities"
	"kadsin/shoot-run/game/helpers"
	"time"
)

func (game *Game) update() {
	ticker := time.NewTicker(time.Millisecond)

	for range ticker.C {
		if game.Exited {
			game.storyShowScore().Show()

			break
		}

		game.generateBlocks()

		game.moveShooter()

		game.generateEnemy()
		game.moveEnemies()

		game.moveBullets()

		game.render()
	}
}

func (game *Game) generateBlocks() {
	if !game.isTimeToGenerateBlocks() {
		return
	}

	game.Blocks = []entities.Object{}

	count := helpers.RandomNumberBetween(10, 15)

	for i := 0; i < count; i++ {
		size := helpers.RandomNumberBetween(5, 10)
		location := helpers.RandomCoordinate(game.Screen, assets.Coordinate{X: 2, Y: 2})

		for j := 0; j < size; j++ {
			isHorizontal := helpers.RandomBoolean()

			shape := '█'
			if isHorizontal {
				shape = '▀'
			}

			block := entities.Object{
				Shape:    shape,
				Location: location,
				Screen:   game.Screen,
				Color:    assets.COLOR_WALLS,
			}

			if isHorizontal {
				location.X++
			} else {
				location.Y++
			}

			game.Blocks = append(game.Blocks, block)
		}
	}
}

func (game *Game) isTimeToGenerateBlocks() bool {
	if time.Now().UnixMilli() > game.LastTimeActions.BlocksGenerator+assets.SPEED_BLOCKS_GENERATOR {
		game.LastTimeActions.BlocksGenerator = time.Now().UnixMilli()
		return true
	}

	return false
}

func (game *Game) moveShooter() {
	if game.isTimeToMoveShooter() {
		if block := game.isShooterBehindOfBlock(); block != nil {
			event := game.EventCollisionBlockByShooter(block)

			if event != nil {
				return
			}
		}

		game.Shooter.Person.UpdateLocation(1)
	}
}

func (game *Game) isTimeToMoveShooter() bool {
	if time.Now().UnixMilli() > game.LastTimeActions.Shooter+int64(game.Shooter.Speed) {
		game.LastTimeActions.Shooter = time.Now().UnixMilli()
		return true
	}

	return false
}

func (game *Game) isShooterBehindOfBlock() *entities.Object {
	for _, block := range game.Blocks {
		if game.Shooter.Person.NextStep(1) == block.Location {
			return &block
		}
	}

	return nil
}

func (game *Game) generateEnemy() {
	if !game.isTimeToGenerateEnemy() {
		return
	}

	enemy := entities.Enemy{
		Person: entities.Object{
			Shape:    '#',
			Location: helpers.RandomCoordinateOnBorders(game.Screen),
			Screen:   game.Screen,
			Color:    assets.COLOR_ENEMIES,
		},
		Target: &game.Shooter.Person,
		Speed:  helpers.RandomNumberBetween(assets.SPEED_MIN_ENEMY, assets.SPEED_MAX_ENEMY),
	}

	game.Enemies = append(game.Enemies, &enemy)
	game.LastTimeActions.Enemies[&enemy] = 0
}

func (game *Game) isTimeToGenerateEnemy() bool {
	if time.Now().UnixMilli() > game.LastTimeActions.EnemyGenerator+int64(game.enemyGeneratorSpeed()) {
		game.LastTimeActions.EnemyGenerator = time.Now().UnixMilli()
		return true
	}

	return false
}

func (game *Game) enemyGeneratorSpeed() uint {
	lastShootDiff := uint(time.Now().UnixMilli()-game.LastTimeActions.Kill) / 100
	if lastShootDiff > assets.SPEED_ENEMY_GENERATOR {
		return assets.SPEED_ENEMY_GENERATOR
	}

	variant := game.KilledEnemiesCount*assets.IMPACT_SHOOT_ON_ENEMY_GENERATING - lastShootDiff
	if variant > 800 {
		return 200
	}

	speed := assets.SPEED_ENEMY_GENERATOR - variant

	return speed
}

func (game *Game) moveEnemies() {
	for _, e := range game.Enemies {
		if !game.isTimeToMoveEnemy(e) {
			continue
		}

		e.LookAtTarget()

		if block := game.isEnemyBehindOfBlock(e); block != nil {
			event := game.EventCollisionBlockByEnemy(block, e)

			if event != nil {
				continue
			}
		}

		e.Person.UpdateLocation(1)

		if e.Person.DoesHit(*e.Target) {
			game.EventCollisionShooterByEnemy(e)
		}
	}
}

func (game *Game) isTimeToMoveEnemy(enemy *entities.Enemy) bool {
	if time.Now().UnixMilli() > game.LastTimeActions.Enemies[enemy]+int64(enemy.Speed) {
		game.LastTimeActions.Enemies[enemy] = time.Now().UnixMilli()
		return true
	}

	return false
}

func (game *Game) isEnemyBehindOfBlock(e *entities.Enemy) *entities.Object {
	for _, block := range game.Blocks {
		if e.Person.NextStep(1) == block.Location {
			return &block
		}
	}

	return nil
}

func (game *Game) moveBullets() {
	if !game.isTimeToMoveBullet() {
		return
	}

	for _, b := range game.Shooter.Bullets {
		if block := game.isBulletBehindOfBlock(b); block != nil {
			game.EventCollisionBlockByBullet(block, b)
		}

		game.Shooter.GoShot(b)

		if enemy := game.anEnemyHitBy(b); enemy != nil {
			game.EventCollisionEnemyByBullet(enemy, b)
		}
	}
}

func (game *Game) isTimeToMoveBullet() bool {
	if time.Now().UnixMilli() > game.LastTimeActions.Bullets+int64(assets.SPEED_BULLET) {
		game.LastTimeActions.Bullets = time.Now().UnixMilli()
		return true
	}

	return false
}

func (game *Game) isBulletBehindOfBlock(bullet *entities.Object) *entities.Object {
	for _, block := range game.Blocks {
		if bullet.DoesHit(block) {
			return &block
		}
	}

	return nil
}

func (game *Game) anEnemyHitBy(bullet *entities.Object) *entities.Enemy {
	for _, e := range game.Enemies {
		if bullet.DoesHit(e.Person) {
			return e
		}
	}

	return nil
}
