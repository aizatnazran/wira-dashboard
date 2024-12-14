<template>
  <div class="combat-type-icon" :title="combatType">
    <img v-if="combatType === 'MELEE'" src="@/assets/combat-types/melee.svg" alt="Melee" />
    <img v-else-if="combatType === 'RANGED'" src="@/assets/combat-types/ranged.svg" alt="Ranged" />
    <img v-else-if="combatType === 'SUPPORT'" src="@/assets/combat-types/support.svg" alt="Support" />
    <img v-else-if="combatType === 'MAGIC'" src="@/assets/combat-types/magic.svg" alt="Magic" />
    <img v-else-if="combatType === 'TANK'" src="@/assets/combat-types/tank.svg" alt="Tank" />
  </div>
</template>

<script>
export default {
  name: 'CombatTypeIcon',
  props: {
    combatType: {
      type: String,
      required: true,
      validator: function(value) {
        return ['MELEE', 'RANGED', 'SUPPORT', 'MAGIC', 'TANK'].includes(value)
      }
    }
  }
}
</script>

<style scoped>
.combat-type-icon {
  position: absolute;
  top: 10px;
  right: 10px;
  width: 32px;
  height: 32px;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 50%;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.combat-type-icon img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  filter: invert(1);
}
</style>
